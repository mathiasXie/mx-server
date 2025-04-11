package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
	"github.com/mathiasXie/gin-web/config"
	"github.com/mathiasXie/gin-web/consts"

	"math/rand"
	"net"
	"strings"
	"time"
)

func UnixtimeToDatetime(timestamp int) string {
	return time.Unix(int64(timestamp), 0).Format(time.DateTime)
}

func DatetimeToUnixtime(datetime string) int64 {
	if t, err := time.ParseInLocation(time.DateTime, datetime, time.Local); err != nil {
		return 0
	} else {
		return t.Unix()
	}
}

func GetUnixTime() int64 {
	return time.Now().Unix()
}

func GetDateTime() string {
	return time.Now().Format(time.DateTime)
}

func GetDate() string {
	return time.Now().Format(time.DateOnly)
}

func LocalIP() string {
	ipList := []string{"114.114.114.114:80", "8.8.8.8:80"}
	for _, ip := range ipList {
		conn, err := net.Dial("udp", ip)
		if err != nil {
			continue
		}
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		conn.Close()
		return localAddr.IP.String()
	}

	return ""
}
func GetClientIP(c *gin.Context) string {
	// 尝试从 X-Forwarded-For 获取
	xffHeader := c.Request.Header.Get("X-Forwarded-For")
	if xffHeader != "" {
		ips := strings.Split(xffHeader, ",")
		return strings.TrimSpace(ips[0])
	}

	// 尝试从 X-Real-IP 获取
	xRealIP := c.Request.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// 最后从 RemoteAddr 获取
	return c.ClientIP()
}

func Hostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func GenerateLogId() string {
	now := time.Now()
	// 使用格式化字符串来指定输出格式，包含毫秒部分
	formattedTime := strings.Replace(now.Format(consts.DateTimeWithoutSpace), ".", "", -1)

	// 初始化随机数生成器，使用当前时间作为种子
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	randomLetters := string(b)
	return fmt.Sprintf("%s%s", formattedTime, randomLetters)
}

const TokenExpireDuration = time.Hour * 24 * 2 // 过期时间 -2天
type TokenClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// MakeToken 生成 jwt token
func MakeToken(userID uint64, username string) (string, error) {
	var claims = TokenClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "mathias",                                  // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.Instance.AccessTokenSecret))
	if err != nil {
		return "", fmt.Errorf("生成token失败:%v", err)
	}
	return signedToken, nil
}

// ParseToken 验证jwt token
func ParseToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (i interface{}, err error) { // 解析token
		return []byte(config.Instance.AccessTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword 验证密码
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
