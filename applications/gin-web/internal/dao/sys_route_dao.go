package dao

import (
	"context"

	"github.com/mathiasXie/curd"
	"github.com/mathiasXie/gin-web/applications/gin-web/internal/model"
	"github.com/mathiasXie/gin-web/pkg/logger"
	"gorm.io/gorm"
)

type SysRouteDAO struct {
	sysRoleRouteModel *curd.Model[model.SysRoleRoute]
	sysRouteModel     *curd.Model[model.SysRoute]
}

type UserRoleRoute struct {
	*model.SysRoleRoute
	Route *model.SysRoute
}

func NewSysRouteDAO(db *gorm.DB) *SysRouteDAO {
	return &SysRouteDAO{
		sysRoleRouteModel: curd.NewModel[model.SysRoleRoute](db),
		sysRouteModel:     curd.NewModel[model.SysRoute](db),
	}
}

func (u *SysRouteDAO) GetRoleRoutes(ctx context.Context, RoleId int64) ([]*UserRoleRoute, error) {
	roleRoutes, _, err := u.sysRoleRouteModel.Select(ctx, &curd.ModelBO{
		Conditions: []map[string]interface{}{
			{"role_id": RoleId},
		},
	})
	if err != nil {
		logger.CtxError(ctx, "[GetRoleRoute]GetRoleRoute Error:", err.Error())
		return nil, err
	}

	if len(roleRoutes) == 0 {
		return []*UserRoleRoute{}, nil
	}

	// 收集所有路由ID
	routeIds := make([]int64, 0, len(roleRoutes))
	for _, v := range roleRoutes {
		routeIds = append(routeIds, v.RouteId)
	}

	// 获取路由信息
	routes, _, err := u.sysRouteModel.Select(ctx, &curd.ModelBO{
		Conditions: []map[string]interface{}{
			{"id": routeIds},
		},
	})
	if err != nil {
		logger.CtxError(ctx, "[GetRoleRoute]GetRoutes Error:", err.Error())
		return nil, err
	}

	// 构建路由映射
	routeMap := make(map[int64]*model.SysRoute)
	for _, route := range routes {
		routeMap[route.Id] = route
	}

	// 组装结果
	result := make([]*UserRoleRoute, 0, len(roleRoutes))
	for _, roleRoute := range roleRoutes {
		if route, ok := routeMap[roleRoute.RouteId]; ok {
			result = append(result, &UserRoleRoute{
				SysRoleRoute: roleRoute,
				Route:        route,
			})
		}
	}

	return result, nil
}

func (u *SysRouteDAO) GetSysRoute(ctx context.Context) ([]*model.SysRoute, error) {
	routes, _, err := u.sysRouteModel.SelectAll(ctx)
	if err != nil {
		logger.CtxError(ctx, "[GetSysRoute]GetSysRoute Error:", err.Error())
		return nil, err
	}
	return routes, nil
}

func (u *SysRouteDAO) BuildMenuRoute(ctx context.Context, userRouteIds []int64, parentId int64) ([]*model.SysRoute, error) {
	routes, _, err := u.sysRouteModel.Select(ctx, &curd.ModelBO{
		Conditions: []map[string]interface{}{
			{"parent_id": parentId, "status": 1, "type": 1, "id": userRouteIds},
		},
		Orders: "sequence asc",
	})

	if err != nil {
		logger.CtxError(ctx, "[BuildMenuRoute]BuildMenuRoute Error:", err.Error())
		return nil, err
	}
	return routes, nil
}

func (u *SysRouteDAO) GetRouteByParentID(ctx context.Context, parentId int64) ([]*model.SysRoute, error) {
	routes, _, err := u.sysRouteModel.Select(ctx, &curd.ModelBO{
		Conditions: []map[string]interface{}{
			{"parent_id": parentId},
		},
	})
	if err != nil {
		logger.CtxError(ctx, "[GetRouteByParentID]GetRouteByParentID Error:", err.Error())
		return nil, err
	}
	return routes, nil
}

func (u *SysRouteDAO) Update(ctx context.Context, Id int64, data map[string]interface{}) error {
	err := u.sysRouteModel.Update(ctx, Id, data)
	if err != nil {
		logger.CtxError(ctx, "[SysRouteDAO]Update Error:", err.Error())
		return err
	}
	return nil
}

func (u *SysRouteDAO) Insert(ctx context.Context, data map[string]interface{}) error {
	_, err := u.sysRouteModel.Insert(ctx, data)
	if err != nil {
		logger.CtxError(ctx, "[SysRouteDAO]Insert Error:", err.Error())
		return err
	}
	return nil
}

func (u *SysRouteDAO) Delete(ctx context.Context, Id int64) error {

	err := u.sysRouteModel.Delete(ctx, &model.SysRoute{Id: Id})
	if err != nil {
		logger.CtxError(ctx, "[SysRouteDAO]Delete Error:", err.Error())
		return err
	}
	return nil
}
