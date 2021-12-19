package calendar

import (
	"context"
	"errors"
	"fmt"
	medCalendar "med-common/app/service/med-calendar/api/v1"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
	"med-common/library/code"

	cdb "git.medlinker.com/service/common/db/v2"
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
)

// 获取项目tree
func (d *Dao) GetList(ctx context.Context, filter map[string]interface{}) (resp []*model.Project, err error) {
	resp = make([]*model.Project, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	if err != nil {
		log.Errorc(ctx, "[dao][FetchAllItem] GetWhereConditionParams(%+v)，err：%+v,err", filter, err)
		return
	}
	err = d.orm.Model(&model.Project{}).Where(whereSql, params...).Find(&resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchAllItem]  GetWhereConditionParams(%+v)，err：%+v,err", filter, err)
	}
	return
}

// 查询项目节点
func (d *Dao) FetchAllItem(ctx context.Context, filter map[string]interface{}) (resp []*model.Project, err error) {
	resp = make([]*model.Project, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	err = d.orm.Model(&model.Project{}).Where(whereSql, params...).Find(&resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchAllItem] 查询失败，err：%+v,err")
	}
	return
}

func (d *Dao) FetchItemPaginate(ctx context.Context, filter map[string]interface{}, page, limit int) (resp []*model.Project, total int64, err error) {
	resp = make([]*model.Project, 0)
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	db := d.orm.Model(&model.Project{}).Where(whereSql, params...)

	err = db.Count(&total).Error
	if err != nil {
		log.Errorc(ctx, "[dao][FetchItemPaginate] 查询count失败，err：%+v,err")
		return
	}
	if total == 0 {
		return
	}
	err = db.Limit(limit).Offset((page - 1) * limit).Find(&resp).Error
	if err = dbutil.IgnoreNoRecordErr(err); err != nil {
		log.Errorc(ctx, "[dao][FetchItemPaginate] 查询失败，err：%+v,err")
	}
	return
}

func (d *Dao) FetchItemCount(ctx context.Context, filter map[string]interface{}) (total int64, err error) {
	db := d.orm.Model(&model.Project{})
	whereSql, params, err := dbutil.GetWhereConditionParams(filter)
	if err = db.Where(whereSql, params...).Count(&total).Error; err != nil {
		log.Errorc(ctx, "[dao][FetchItemCount] 查询count失败，err：%+v,err")
		return
	}

	return
}

func (d *Dao) UpdateItem(ctx context.Context, project *model.Project, req *medCalendar.EditItemReq) (err error) {
	// 父项目不能是自身
	log.Infoc(ctx, "UpdateItem,project:%+v,req:%+v", project, req)
	if req.ParentId == project.Id {
		return errors.New("父项目不能设置为自身")
	}
	// 检查父项目编码是否发生了变更，如果是,则需要重新,并更新编码

	projectMap := make(map[string]interface{})
	projectMap["parent_id"] = req.ParentId
	projectMap["project_name"] = req.ProjectName
	projectMap["cost_dept_code"] = req.CostDeptCode
	projectMap["remark"] = req.Remark
	projectMap["is_required_comment"] = req.IsRequiredComment
	projectMap["is_force_notice"] = req.IsForceNotice
	projectMap["notice"] = req.Notice

	updateTask := func(db *gorm.DB) error {
		// 更新project
		err := db.Model(&model.Project{}).Where("id = ?", req.ProjectId).Updates(projectMap).Error
		if err != nil {
			return fmt.Errorf("更新项目失败 <err:%v>", err)
		}

		return nil
	}

	if err = cdb.DoTrans(d.orm, updateTask); err != nil {
		log.Errorc(ctx, "更新失败 <err:%v> <req:%+v>", err, req)
		err = code.MedErrDBFailed
		return
	}
	// 如果成本部门发生变更，需要更新project_task
	if req.CostDeptCode != "" && project.CostDeptCode != req.CostDeptCode {
		err = d.orm.Exec("update project_task set cost_dept_code = ? where project_id =?",
			req.CostDeptCode, req.ProjectId).Error
		if err != nil {
			return err
		}
	}
	return
}

func (d *Dao) DeleteItem(ctx context.Context, req *medCalendar.DeleteItemReq) (err error) {
	err = d.orm.Model(&model.Project{}).
		Where("id = ?", req.ProjectId).
		UpdateColumn("state", -1).Error
	return
}

// 创建
func (d *Dao) CreateItem(ctx context.Context, project *model.Project) (err error) {
	db := d.orm.Model(&model.Project{})
	err = db.Create(project).Error
	if err != nil {
		log.Errorc(ctx, "[dao][CreateItem] 创建日常事务，err：%+v,", err)
	}
	return
}

func (d *Dao) GetItemByCodes(ctx context.Context, ids []int64) (itemMap map[int64]*model.Project, err error) {
	itemMap = make(map[int64]*model.Project, 0)
	var items []*model.Project
	err = d.orm.Where("id in (?)", ids).Find(&items).Error
	for _, i := range items {
		itemMap[i.Id] = i
	}
	return
}
