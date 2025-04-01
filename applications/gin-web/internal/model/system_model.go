// Code generated from sys.sql by gen_sql_model. DO NOT EDIT.

package model

import (
	"time"
)

func (SysUser) TableName() string {
	return "sys_users"
}

type SysUser struct {
	IsSuper   int8      `gorm:"column:is_super" json:"is_super"` // 是否为超级管理员 1:是 0:否
	Status    int8      `gorm:"column:status" json:"status"`     // 状态：1：正常; 0：禁用
	Id        int32     `gorm:"column:id;primary_key" json:"id"`
	RoleId    int32     `gorm:"column:role_id" json:"role_id"`   // 所属权限组ID
	Email     string    `gorm:"column:email" json:"email"`       // 邮箱
	LastIp    string    `gorm:"column:last_ip" json:"last_ip"`   // 最后一次操作IP
	Mobile    string    `gorm:"column:mobile" json:"mobile"`     // 手机号码
	Password  string    `gorm:"column:password" json:"password"` // 密码
	Realname  string    `gorm:"column:realname" json:"realname"` // 真实姓名
	Username  string    `gorm:"column:username" json:"username"` // 用户名
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	LastTime  time.Time `gorm:"column:last_time" json:"last_time"` // 最后一次操作时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewSysUser() *SysUser {
	return &SysUser{}
}

func (SysRole) TableName() string {
	return "sys_roles"
}

type SysRole struct {
	Status      int8      `gorm:"column:status" json:"status"`     // 状态 1=启用;0=禁用
	Sequence    int32     `gorm:"column:sequence" json:"sequence"` // 排序
	Id          int64     `gorm:"column:id;primary_key" json:"id"`
	ParentId    int64     `gorm:"column:parent_id" json:"parent_id"`     // 上级角色ID
	Description string    `gorm:"column:description" json:"description"` // 角色描述
	Name        string    `gorm:"column:name" json:"name"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewSysRole() *SysRole {
	return &SysRole{}
}

func (SysRoleRoute) TableName() string {
	return "sys_role_routes"
}

type SysRoleRoute struct {
	Id      int64 `gorm:"column:id;primary_key" json:"id"`
	RoleId  int64 `gorm:"column:role_id" json:"role_id"`   // 角色id
	RouteId int64 `gorm:"column:route_id" json:"route_id"` // 路由功能表
}

func NewSysRoleRoute() *SysRoleRoute {
	return &SysRoleRoute{}
}

func (SysRoute) TableName() string {
	return "sys_routes"
}

type SysRoute struct {
	Display    int8      `gorm:"column:display" json:"display"`       // 0=隐藏;1=显示
	Global     int8      `gorm:"column:global" json:"global"`         // 公共资源 1是,2否 无需分配所有人就可以访问的
	Status     int8      `gorm:"column:status" json:"status"`         // 1=启用;0=禁用
	Type       int8      `gorm:"column:type" json:"type"`             // 类型（1=菜单;2=按钮）
	Sequence   int32     `gorm:"column:sequence" json:"sequence"`     // 排序
	Id         int64     `gorm:"column:id" json:"id"`                 // 主键
	ParentId   int64     `gorm:"column:parent_id" json:"parent_id"`   // 父级菜单ID
	ApiPath    string    `gorm:"column:api_path" json:"api_path"`     // api地址,不包含management,用于中间件鉴权,如果此功能使用到多个api用逗号隔开
	Component  string    `gorm:"column:component" json:"component"`   // 组件
	Icon       string    `gorm:"column:icon" json:"icon"`             // 菜单图标
	Label      string    `gorm:"column:label" json:"label"`           // 名称
	Name       string    `gorm:"column:name" json:"name"`             // 路由标识
	Path       string    `gorm:"column:path" json:"path"`             // 路径;对应前端vue路由
	Permission string    `gorm:"column:permission" json:"permission"` // 权限标识,用于前端按钮级别的鉴权
	Style      string    `gorm:"column:style" json:"style"`           // 样式
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func NewSysRoute() *SysRoute {
	return &SysRoute{}
}
