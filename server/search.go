package calendar

import (
	"context"
	"fmt"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
	"strings"

	"github.com/go-kratos/kratos/pkg/sync/errgroup"
	"gorm.io/gorm"
)

func (d *Dao) GetProjectData(where map[string]interface{}, order ...string) (items []*model.Project, err error) {
	if len(where) == 0 {
		return
	}
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	db := d.orm.Model(&model.Project{}).Where(whereSql, params...)

	orders := strings.Join(order, ",")
	if orders != "" {
		db = db.Order(orders)
	}
	err = db.Find(&items).Error
	if err != nil {
		return
	}
	return
}

// 获取查询结果对应项目的1级项目离诶包数据
func (d *Dao) GetProjectLevelOneData(where map[string]interface{}, order ...string) (items []*model.Project, err error) {
	if len(where) == 0 {
		return
	}
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	db := d.orm.Model(&model.Project{}).Where(whereSql, params...)

	orders := strings.Join(order, ",")
	if orders != "" {
		db = db.Order(orders)
	}
	err = db.Find(&items).Error
	if err != nil {
		return
	}
	if pid, ok := where["parent_id"]; ok && pid == 0 {
		return
	}
	fmt.Println("-----", len(items))
	// 如果有非一级项目将其上级项目显示出来
	levelOneProjectIds := make([]int64, 0)
	for _, v := range items {
		if v.ParentId > 0 {
			levelOneProjectIds = append(levelOneProjectIds, v.ParentId)
			continue
		}
		levelOneProjectIds = append(levelOneProjectIds, v.Id)
	}
	err = d.orm.Model(&model.Project{}).Where("id in (?)", levelOneProjectIds).Order(orders).Find(&items).Error
	if err != nil {
		return
	}

	return
}

func (d *Dao) GetProjectRiskData() (m map[int64][]*model.ProjectRisk, err error) {
	m = make(map[int64][]*model.ProjectRisk)
	var items []*model.ProjectRisk
	err = d.orm.Model(&model.ProjectRisk{}).Find(&items).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	if err != nil {
		return
	}
	for _, item := range items {
		m[item.ProjectId] = append(m[item.ProjectId], item)
	}
	return
}

func (d *Dao) ProjectNodeCount() (ret map[int64]int64, err error) {
	type NodeCount struct {
		ProjectId int64
		Cnt       int64
	}

	var nc []*NodeCount
	ret = make(map[int64]int64, 0)
	err = d.orm.Model(&model.ProjectNode{}).
		Select("project_id,count(*) as cnt").
		Where("is_auto_create = ? ", 2).
		Group("project_id").
		Scan(&nc).Error

	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	for _, n := range nc {
		ret[n.ProjectId] = n.Cnt
	}
	return
}

func (d *Dao) ChildrenProjectCount() (ret map[int64]int64, err error) {
	type ProjectCount struct {
		ParentId int64
		Cnt      int64
	}

	var pc []*ProjectCount
	ret = make(map[int64]int64, 0)
	err = d.orm.Model(&model.Project{}).
		Where("state != ?", -1).
		Select("parent_id,count(*) as cnt").
		Group("parent_id").
		Scan(&pc).Error

	if err == gorm.ErrRecordNotFound {
		err = nil
	}

	for _, p := range pc {
		ret[p.ParentId] = p.Cnt
	}
	return
}

func (d *Dao) GetProjectNodeData(where map[string]interface{}) (items []*model.ProjectNode, err error) {
	whereSql, params, err := dbutil.GetWhereConditionParams(where)
	if err != nil {
		return
	}
	err = d.orm.Model(&model.ProjectNode{}).Where(whereSql, params...).Find(&items).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	if err != nil {
		return
	}
	return
}

func (d *Dao) GetProjectNodeMap() (m map[int64][]*model.ProjectNode, err error) {
	m = make(map[int64][]*model.ProjectNode)
	var items []*model.ProjectNode
	err = d.orm.Model(&model.ProjectNode{}).Find(&items).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	if err != nil {
		return
	}
	for _, item := range items {
		m[item.ProjectId] = append(m[item.ProjectId], item)
	}
	return
}

func (d *Dao) GetProjectNodeItemData() (m map[int64][]*model.ProjectNodeItem, err error) {
	m = make(map[int64][]*model.ProjectNodeItem)
	var items []*model.ProjectNodeItem
	err = d.orm.Model(&model.ProjectNodeItem{}).Find(&items).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	if err != nil {
		return
	}
	for _, item := range items {
		m[item.NodeId] = append(m[item.NodeId], item)
	}
	return
}

type ProjectAuthData struct {
	Id       int64
	AuthDep  map[int64]*model.ProjectAuthDep
	AuthUser map[string]*model.ProjectAuthUser
}

func (d *Dao) GetProjectAuthData(ctx context.Context) (proMap map[int64]*ProjectAuthData, err error) {
	var authDep []*model.ProjectAuthDep
	var authUser []*model.ProjectAuthUser
	eg := errgroup.WithContext(ctx)
	eg.Go(func(ctx context.Context) error {
		err = d.orm.Model(&model.ProjectAuthDep{}).Find(&authDep).Error
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return err
	})
	eg.Go(func(ctx context.Context) error {
		err = d.orm.Model(&model.ProjectAuthUser{}).Find(&authUser).Error
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return err
	})
	err = eg.Wait()
	proMap = make(map[int64]*ProjectAuthData)
	for _, item := range authDep {
		if _, ok := proMap[item.ProjectId]; !ok {
			proMap[item.ProjectId] = &ProjectAuthData{
				Id:       item.ProjectId,
				AuthDep:  map[int64]*model.ProjectAuthDep{},
				AuthUser: map[string]*model.ProjectAuthUser{},
			}
		}
		proMap[item.ProjectId].AuthDep[item.DeptId] = item
	}
	for _, item := range authUser {
		if _, ok := proMap[item.ProjectId]; !ok {
			proMap[item.ProjectId] = &ProjectAuthData{
				Id:       item.ProjectId,
				AuthDep:  map[int64]*model.ProjectAuthDep{},
				AuthUser: map[string]*model.ProjectAuthUser{},
			}
		}
		proMap[item.ProjectId].AuthUser[item.UserNo] = item
	}
	return
}

func (d *Dao) GetProjectFocusUsers() (users map[int64][]*model.User, err error) {
	var focus []*model.UserFocusProject

	users = make(map[int64][]*model.User, 0)
	err = d.orm.Table("user_focus_project").
		Select("distinct user_focus_project.user_id,user_focus_project.project_id").
		Joins(" left join user on user_focus_project.user_id = user.user_id ").
		Scan(&focus).Error

	if len(focus) == 0 {
		err = nil
		return
	}

	// 获取用户名称
	userMap, err := d.FetchAllUserMap()
	if err != nil {
		return
	}

	for _, f := range focus {
		if user, ok := userMap[f.UserId]; ok {
			users[f.ProjectId] = append(users[f.ProjectId], user)
		}
	}
	return
}
