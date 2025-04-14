package utils

/*
#cgo CFLAGS: -I/root/miniconda3/include
#cgo LDFLAGS: -L/root/miniconda3/lib -lavformat -lavcodec -lavutil -lswresample -lopus

#include <libavformat/avformat.h>
#include <libavcodec/avcodec.h>
#include <libavutil/opt.h>
#include <libavutil/samplefmt.h>
#include <libswresample/swresample.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// 初始化输入和输出格式上下文、编解码器上下文等
int init_conversion(const char *input_file, AVFormatContext **input_format_context,
                    AVCodecContext **input_codec_context,
                    AVFormatContext **output_format_context,
                    AVCodecContext **output_codec_context,
                    int *audio_stream_index) {
    int ret;

    // 打开输入文件
    if ((ret = avformat_open_input(input_format_context, input_file, NULL, NULL)) < 0) {
        fprintf(stderr, "Could not open input file\n");
        return ret;
    }

    // 读取文件的流信息
    if ((ret = avformat_find_stream_info(*input_format_context, NULL)) < 0) {
        fprintf(stderr, "Could not find stream information\n");
        avformat_close_input(input_format_context);
        return ret;
    }

    // 查找音频流
    *audio_stream_index = av_find_best_stream(*input_format_context, AVMEDIA_TYPE_AUDIO, -1, -1, NULL, 0);
    if (*audio_stream_index < 0) {
        fprintf(stderr, "Could not find audio stream\n");
        avformat_close_input(input_format_context);
        return *audio_stream_index;
    }

    // 查找输入编解码器
    AVCodec *input_codec = avcodec_find_decoder((*input_format_context)->streams[*audio_stream_index]->codecpar->codec_id);
    if (!input_codec) {
        fprintf(stderr, "Could not find input codec\n");
        avformat_close_input(input_format_context);
        return AVERROR(EINVAL);
    }

    // 分配输入编解码器上下文
    *input_codec_context = avcodec_alloc_context3(input_codec);
    if (!*input_codec_context) {
        fprintf(stderr, "Could not allocate input codec context\n");
        avformat_close_input(input_format_context);
        return AVERROR(ENOMEM);
    }

    // 填充输入编解码器上下文
    if ((ret = avcodec_parameters_to_context(*input_codec_context, (*input_format_context)->streams[*audio_stream_index]->codecpar)) < 0) {
        fprintf(stderr, "Could not copy input codec parameters to context\n");
        avcodec_free_context(input_codec_context);
        avformat_close_input(input_format_context);
        return ret;
    }

    // 打开输入编解码器
    if ((ret = avcodec_open2(*input_codec_context, input_codec, NULL)) < 0) {
        fprintf(stderr, "Could not open input codec\n");
        avcodec_free_context(input_codec_context);
        avformat_close_input(input_format_context);
        return ret;
    }

    // 查找 Opus 编码器
    AVCodec *output_codec = avcodec_find_encoder(AV_CODEC_ID_OPUS);
    if (!output_codec) {
        fprintf(stderr, "Could not find Opus encoder\n");
        avcodec_free_context(input_codec_context);
        avformat_close_input(input_format_context);
        return AVERROR(EINVAL);
    }

    // 分配输出编解码器上下文
    *output_codec_context = avcodec_alloc_context3(output_codec);
    if (!*output_codec_context) {
        fprintf(stderr, "Could not allocate output codec context\n");
        avcodec_free_context(input_codec_context);
        avformat_close_input(input_format_context);
        return AVERROR(ENOMEM);
    }

    // 设置输出编解码器参数
    (*output_codec_context)->sample_fmt = AV_SAMPLE_FMT_S16;
    (*output_codec_context)->sample_rate = 16000;
    (*output_codec_context)->channel_layout = AV_CH_LAYOUT_MONO;
    (*output_codec_context)->channels = 1;

    // 打开输出编解码器
    if ((ret = avcodec_open2(*output_codec_context, output_codec, NULL)) < 0) {
        fprintf(stderr, "Could not open output codec\n");
        avcodec_free_context(input_codec_context);
        avcodec_free_context(output_codec_context);
        avformat_close_input(input_format_context);
        return ret;
    }

    // 分配输出格式上下文
    avformat_alloc_output_context2(output_format_context, NULL, "opus", NULL);
    if (!*output_format_context) {
        fprintf(stderr, "Could not allocate output format context\n");
        avcodec_free_context(input_codec_context);
        avcodec_free_context(output_codec_context);
        avformat_close_input(input_format_context);
        return AVERROR(ENOMEM);
    }

    // 创建输出流
    AVStream *output_stream = avformat_new_stream(*output_format_context, output_codec);
    if (!output_stream) {
        fprintf(stderr, "Could not create output stream\n");
        avcodec_free_context(input_codec_context);
        avcodec_free_context(output_codec_context);
        avformat_free_context(*output_format_context);
        avformat_close_input(input_format_context);
        return AVERROR(ENOMEM);
    }

    // 复制编解码器参数到输出流
    if ((ret = avcodec_parameters_from_context(output_stream->codecpar, *output_codec_context)) < 0) {
        fprintf(stderr, "Could not copy codec parameters to output stream\n");
        avcodec_free_context(input_codec_context);
        avcodec_free_context(output_codec_context);
        avformat_free_context(*output_format_context);
        avformat_close_input(input_format_context);
        return ret;
    }

    return 0;
}

// 转换音频为 Opus 编码
int convert_audio(AVFormatContext *input_format_context, AVCodecContext *input_codec_context,
                  AVFormatContext *output_format_context, AVCodecContext *output_codec_context,
                  int audio_stream_index, unsigned char ***opus_datas, int *opus_data_count,
                  double *duration) {
    int ret;
    AVPacket *input_packet = av_packet_alloc();
    AVFrame *input_frame = av_frame_alloc();
    AVFrame *output_frame = av_frame_alloc();
    AVPacket *output_packet = av_packet_alloc();

    if (!input_packet || !input_frame || !output_frame || !output_packet) {
        fprintf(stderr, "Could not allocate packet or frame\n");
        ret = AVERROR(ENOMEM);
        goto end;
    }

    // 初始化重采样上下文
    SwrContext *swr_context = swr_alloc();
    if (!swr_context) {
        fprintf(stderr, "Could not allocate resampler context\n");
        ret = AVERROR(ENOMEM);
        goto end;
    }

    // 设置重采样参数
    av_opt_set_int(swr_context, "in_channel_layout", input_codec_context->channel_layout, 0);
    av_opt_set_int(swr_context, "out_channel_layout", output_codec_context->channel_layout, 0);
    av_opt_set_int(swr_context, "in_sample_rate", input_codec_context->sample_rate, 0);
    av_opt_set_int(swr_context, "out_sample_rate", output_codec_context->sample_rate, 0);
    av_opt_set_sample_fmt(swr_context, "in_sample_fmt", input_codec_context->sample_fmt, 0);
    av_opt_set_sample_fmt(swr_context, "out_sample_fmt", output_codec_context->sample_fmt, 0);

    // 初始化重采样上下文
    if ((ret = swr_init(swr_context)) < 0) {
        fprintf(stderr, "Could not initialize resampler context\n");
        goto end;
    }

    // 计算音频时长
    *duration = (double)input_format_context->streams[audio_stream_index]->duration *
                av_q2d(input_format_context->streams[audio_stream_index]->time_base);

    *opus_datas = NULL;
    *opus_data_count = 0;

    // 读取输入包
    while (av_read_frame(input_format_context, input_packet) >= 0) {
        if (input_packet->stream_index == audio_stream_index) {
            // 发送输入包到解码器
            if ((ret = avcodec_send_packet(input_codec_context, input_packet)) < 0) {
                fprintf(stderr, "Error sending packet to decoder\n");
                goto end;
            }

            // 从解码器接收帧
            while (ret >= 0) {
                ret = avcodec_receive_frame(input_codec_context, input_frame);
                if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
                    break;
                } else if (ret < 0) {
                    fprintf(stderr, "Error receiving frame from decoder\n");
                    goto end;
                }

                // 重采样帧
                av_frame_copy_props(output_frame, input_frame);
                output_frame->format = output_codec_context->sample_fmt;
                output_frame->channel_layout = output_codec_context->channel_layout;
                output_frame->sample_rate = output_codec_context->sample_rate;
                output_frame->nb_samples = av_rescale_rnd(input_frame->nb_samples,
                                                          output_codec_context->sample_rate,
                                                          input_codec_context->sample_rate,
                                                          AV_ROUND_UP);

                if ((ret = av_frame_get_buffer(output_frame, 0)) < 0) {
                    fprintf(stderr, "Could not allocate output frame buffer\n");
                    goto end;
                }

                ret = swr_convert(swr_context, output_frame->data, output_frame->nb_samples,
                                  (const uint8_t **)input_frame->data, input_frame->nb_samples);
                if (ret < 0) {
                    fprintf(stderr, "Error resampling audio\n");
                    goto end;
                }

                // 发送重采样后的帧到编码器
                if ((ret = avcodec_send_frame(output_codec_context, output_frame)) < 0) {
                    fprintf(stderr, "Error sending frame to encoder\n");
                    goto end;
                }

                // 从编码器接收编码后的包
                while (ret >= 0) {
                    ret = avcodec_receive_packet(output_codec_context, output_packet);
                    if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
                        break;
                    } else if (ret < 0) {
                        fprintf(stderr, "Error receiving packet from encoder\n");
                        goto end;
                    }

                    // 保存编码后的包数据
                    *opus_datas = realloc(*opus_datas, (*opus_data_count + 1) * sizeof(unsigned char *));
                    (*opus_datas)[*opus_data_count] = malloc(output_packet->size);
                    memcpy((*opus_datas)[*opus_data_count], output_packet->data, output_packet->size);
                    (*opus_data_count)++;

                    av_packet_unref(output_packet);
                }
            }
        }
        av_packet_unref(input_packet);
    }

    // 刷新编码器
    ret = avcodec_send_frame(output_codec_context, NULL);
    while (ret >= 0) {
        ret = avcodec_receive_packet(output_codec_context, output_packet);
        if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
            break;
        } else if (ret < 0) {
            fprintf(stderr, "Error receiving packet from encoder during flush\n");
            goto end;
        }

        // 保存编码后的包数据
        *opus_datas = realloc(*opus_datas, (*opus_data_count + 1) * sizeof(unsigned char *));
        (*opus_datas)[*opus_data_count] = malloc(output_packet->size);
        memcpy((*opus_datas)[*opus_data_count], output_packet->data, output_packet->size);
        (*opus_data_count)++;

        av_packet_unref(output_packet);
    }

    ret = 0;

end:
    av_packet_free(&input_packet);
    av_frame_free(&input_frame);
    av_frame_free(&output_frame);
    av_packet_free(&output_packet);
    swr_free(&swr_context);

    return ret;
}

// 清理资源
void cleanup(AVFormatContext *input_format_context, AVCodecContext *input_codec_context,
             AVFormatContext *output_format_context, AVCodecContext *output_codec_context,
             unsigned char **opus_datas, int opus_data_count) {
    for (int i = 0; i < opus_data_count; i++) {
        free(opus_datas[i]);
    }
    free(opus_datas);
    avcodec_free_context(&input_codec_context);
    avformat_close_input(&input_format_context);
    avcodec_free_context(&output_codec_context);
    avformat_free_context(output_format_context);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// AudioToOpusData 将音频文件转换为 Opus 编码数据
func AudioToOpusData(audioFilePath string) ([][]byte, float64) {
	var (
		inputFormatContext  *C.AVFormatContext
		inputCodecContext   *C.AVCodecContext
		outputFormatContext *C.AVFormatContext
		outputCodecContext  *C.AVCodecContext
		audioStreamIndex    C.int
		opusDatas           **C.uchar
		opusDataCount       C.int
		duration            C.double
	)

	audioFilePathC := C.CString(audioFilePath)
	defer C.free(unsafe.Pointer(audioFilePathC))

	// 初始化转换
	if ret := C.init_conversion(audioFilePathC, &inputFormatContext, &inputCodecContext,
		&outputFormatContext, &outputCodecContext, &audioStreamIndex); ret != 0 {
		fmt.Printf("Initialization failed: %d\n", int(ret))
		return nil, 0
	}

	// 转换音频
	if ret := C.convert_audio(inputFormatContext, inputCodecContext,
		outputFormatContext, outputCodecContext, audioStreamIndex,
		&opusDatas, &opusDataCount, &duration); ret != 0 {
		fmt.Printf("Conversion failed: %d\n", int(ret))
		C.cleanup(inputFormatContext, inputCodecContext, outputFormatContext, outputCodecContext, opusDatas, opusDataCount)
		return nil, 0
	}

	// 将 C 数组转换为 Go 切片
	opusDataSlices := make([][]byte, int(opusDataCount))
	for i := 0; i < int(opusDataCount); i++ {
		var packetSize int
		for j := 0; ; j++ {
			if *(*C.uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(opusDatas)) + uintptr(i)*unsafe.Sizeof(*opusDatas) + uintptr(j))) == 0 {
				packetSize = j
				break
			}
		}
		opusDataSlices[i] = make([]byte, packetSize)
		for j := 0; j < packetSize; j++ {
			opusDataSlices[i][j] = byte(*(*C.uchar)(unsafe.Pointer(uintptr(unsafe.Pointer(opusDatas)) + uintptr(i)*unsafe.Sizeof(*opusDatas) + uintptr(j))))
		}
	}

	// 清理资源
	C.cleanup(inputFormatContext, inputCodecContext, outputFormatContext, outputCodecContext, opusDatas, opusDataCount)

	return opusDataSlices, float64(duration)
}
