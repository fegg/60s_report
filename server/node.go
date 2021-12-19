package calendar

import (
	"context"
	medCalendar "med-common/app/service/med-calendar/api/v1"
	"med-common/app/service/med-calendar/model"
	"med-common/app/service/prescription-service/pkg/dbutil"
	"med-common/library/code"
	"time"

	cdb "git.medlinker.com/service/common/db/v2"
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
)

// 查询项目节点
func (d *Dao) OneNode(nodeId int64) (node *model.ProjectNode, err error) {
	node = &model.ProjectNode{}
	err = d.orm.Model(&model.ProjectNode{}).First(node, "id = ?", nodeId).Error
	return
}

// 添加项目节点
func (d *Dao) AddProjectNode(ctx context.Context, node *model.ProjectNode) (err error) {
	if err = d.orm.Create(node).Error; err != nil {
		log.Errorc(ctx, "AddProjectNode err <params:%+v> <err:%v>", node, err)
		err = code.MedErrDBFailed
		return
	}
	return
}

// 添加项目节点
func (d *Dao) BatchAddProjectNode(ctx context.Context, nodes []*model.ProjectNode) (err error) {
	if err = d.orm.Create(nodes).Error; err != nil {
		log.Errorc(ctx, "AddProjectNode err <params:%+v> <err:%v>", nodes, err)
		err = code.MedErrDBFailed
		return
	}
	return
}

// 获取待审核节点列表
func (d *Dao) WaitAuditNodeList(ctx context.Context, req *medCalendar.GetNodeListReq) (nodes []model.ProjectNode, count int64, err error) {
	node := model.ProjectNode{}
	// 查询待审批
	db := d.orm.Table(node.TableName())
	if len(req.Status) > 0 {
		db = db.Where("status in (?)", req.Status)
	}
	if req.ProjectId != 0 {
		db = db.Where("project_id = ?", req.ProjectId)
	}
	db = db.Where("is_auto_create = ?", 2)
	db = db.Where("state = 1")
	if err = db.Count(&count).Error; err != nil {
		log.Errorc(ctx, "WaitAuditNodeList.dao 获取节点列表数量失败 <err:%v>", err)
		return
	}
	db = db.Offset(int(req.Start)).Limit(int(req.Limit))
	if len(req.Sorts) == 0 {
		req.Sorts = []string{"finish_time"}
	}
	db = db.Order(dbutil.BuildOrder("project_node", req.Sorts))
	if err = db.Find(&nodes).Error; nil != err {
		log.Errorc(ctx, "WaitAuditNodeList.dao 获取节点列表失败 <err:%v>", err.Error())
		return
	}
	return
}

// 审核项目节点
func (d *Dao) AuditProjectNode(ctx context.Context, node *model.ProjectNode, req *medCalendar.AuditProjectNodeReq) (err error) {
	beforeStatus := node.Status
	var afterStatus int64
	maps := make(map[string]interface{})
	if req.IsPass {
		afterStatus = 3
	} else {
		afterStatus = 4
	}
	maps["status"] = afterStatus
	maps["audit_time"] = time.Now().Unix()
	auditLog := &model.AuditNodeLog{
		NodeId:        int64(node.Id),
		AuditUserId:   req.UserId,
		AuditUserName: req.UserName,
		Bofore:        beforeStatus,
		After:         afterStatus,
		Remark:        req.Remark,
	}
	auditTask := func(db *gorm.DB) error {
		err := db.Model(&model.ProjectNode{}).Where("id = ?", node.Id).Where("status = 1").Updates(maps).Error
		if err != nil {
			return err
		}
		err = db.Model(&model.AuditNodeLog{}).Create(auditLog).Error
		if err != nil {
			return err
		}
		return nil
	}
	if err = cdb.DoTrans(d.orm, auditTask); err != nil {
		log.Errorc(ctx, "审核项目节点失败 <err:%v> <req:%+v>", err, req)
		return
	}
	return
}

// 删除节点
func (d *Dao) DeleteProjectNode(ctx context.Context, node *model.ProjectNode) (err error) {
	return d.orm.Model(&model.ProjectNode{}).Where("id = ?", node.Id).Delete(node).Error
}

// 更新项目节点
func (d *Dao) ModifyProjectNode(ctx context.Context, nodeId int64, node map[string]interface{}) (err error) {
	return d.orm.Model(&model.ProjectNode{}).Where("id = ?", nodeId).Updates(node).Error
}
