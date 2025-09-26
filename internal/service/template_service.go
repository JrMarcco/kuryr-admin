package service

import (
	"context"
	"fmt"
	"time"

	"github.com/JrMarcco/easy-grpc/client"
	"github.com/JrMarcco/kuryr-admin/internal/domain"
	commonv1 "github.com/JrMarcco/kuryr-api/api/go/common/v1"
	templatev1 "github.com/JrMarcco/kuryr-api/api/go/template/v1"
)

type TemplateService interface {
	Save(ctx context.Context, tpl domain.ChannelTemplate) (domain.ChannelTemplate, error)

	SaveVersion(ctx context.Context, version domain.ChannelTemplateVersion) (domain.ChannelTemplateVersion, error)
}

var _ TemplateService = (*DefaultTemplateService)(nil)

type DefaultTemplateService struct {
	grpcServerName string
	grpcClients    *client.Manager[templatev1.TemplateServiceClient]
}

func (s *DefaultTemplateService) Save(ctx context.Context, tpl domain.ChannelTemplate) (domain.ChannelTemplate, error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return domain.ChannelTemplate{}, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := grpcClient.SaveTemplate(ctx, &templatev1.SaveTemplateRequest{Template: s.domainToPb(tpl)})
	if err != nil {
		return domain.ChannelTemplate{}, fmt.Errorf("[kuryr-admin] failed to save template: %w", err)
	}
	return s.pbToDomain(resp.Template), nil
}

func (s *DefaultTemplateService) SaveVersion(ctx context.Context, version domain.ChannelTemplateVersion) (domain.ChannelTemplateVersion, error) {
	grpcClient, err := s.grpcClients.Get(s.grpcServerName)
	if err != nil {
		return domain.ChannelTemplateVersion{}, fmt.Errorf("[kuryr-admin] failed to get grpc client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := grpcClient.SaveTemplateVersion(
		ctx,
		&templatev1.SaveTemplateVersionRequest{Version: s.domainToPbVersion(version)},
	)
	if err != nil {
		return domain.ChannelTemplateVersion{}, fmt.Errorf("[kuryr-admin] failed to save template version: %w", err)
	}
	return s.pbToDomainVersion(resp.Version), nil
}

func (s *DefaultTemplateService) domainToPb(tpl domain.ChannelTemplate) *templatev1.ChannelTemplate {
	return &templatev1.ChannelTemplate{
		BizId:            tpl.BizId,
		BizType:          string(tpl.BizType),
		TplName:          tpl.TplName,
		TplDesc:          tpl.TplDesc,
		Channel:          commonv1.Channel(tpl.Channel),
		NotificationType: tpl.NotificationType,
	}
}

func (s *DefaultTemplateService) domainToPbVersion(version domain.ChannelTemplateVersion) *templatev1.TemplateVersion {
	return &templatev1.TemplateVersion{
		TplId:       version.TplId,
		VersionName: version.VersionName,
		Signature:   version.Signature,
		Content:     version.Content,
		ApplyRemark: version.ApplyRemark,
	}
}

func (s *DefaultTemplateService) pbToDomain(pb *templatev1.ChannelTemplate) domain.ChannelTemplate {
	return domain.ChannelTemplate{
		Id:                 pb.Id,
		BizId:              pb.BizId,
		BizType:            domain.BizType(pb.BizType),
		TplName:            pb.TplName,
		TplDesc:            pb.TplDesc,
		Channel:            int32(pb.Channel),
		NotificationType:   int32(pb.NotificationType),
		ActivatedVersionId: pb.ActivatedVersionId,
		CreatedAt:          pb.CreatedAt,
		UpdatedAt:          pb.UpdatedAt,
	}
}

func (s *DefaultTemplateService) pbToDomainVersion(pb *templatev1.TemplateVersion) domain.ChannelTemplateVersion {
	return domain.ChannelTemplateVersion{
		Id:              pb.Id,
		TplId:           pb.TplId,
		VersionName:     pb.VersionName,
		Signature:       pb.Signature,
		Content:         pb.Content,
		ApplyRemark:     pb.ApplyRemark,
		AuditId:         pb.AuditId,
		AuditorId:       pb.AuditorId,
		AuditTime:       pb.AuditTime,
		AuditStatus:     pb.AuditStatus,
		RejectionReason: pb.RejectionReason,
		LastReviewAt:    pb.LastReviewAt,
		CreatedAt:       pb.CreatedAt,
		UpdatedAt:       pb.UpdatedAt,
	}
}

func NewDefaultTemplateService(
	grpcServerName string,
	grpcClients *client.Manager[templatev1.TemplateServiceClient],
) *DefaultTemplateService {
	return &DefaultTemplateService{
		grpcServerName: grpcServerName,
		grpcClients:    grpcClients,
	}
}
