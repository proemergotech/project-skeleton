package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.com/proemergotech/geb-client-go/v2/geb"
	"gitlab.com/proemergotech/log-go/v3"
	"gitlab.com/proemergotech/shallow-go"
	"gitlab.com/proemergotech/uuid-go"
	yafuds "gitlab.com/proemergotech/yafuds-client-go/v2/client"

	//%:{{ `
	"gitlab.com/proemergotech/dliver-project-skeleton/app/client"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/centrifuge"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/skeleton"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/validation"
	//%: ` | replace "dliver-project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
)

type Service struct {
	//%: {{- if .Centrifuge }}
	centrifugeClient *client.Centrifuge
	centrifugeJSON   jsoniter.API
	//%: {{- end }}
	//%: {{- if .Yafuds }}
	yafudsClient yafuds.Client
	validator    *validation.Validator
	//%: {{- end }}
	//%: {{- if and .Geb .Examples }}
	gebQueue *geb.Queue
	//%: {{- end }}
	//%: {{- if .SiteConfig }}
	siteConfigClient *client.SiteConfig
	//%: {{- end }}
}

func NewService(
	//%: {{- if .Centrifuge }}
	centrifugeClient *client.Centrifuge,
	centrifugeJSON jsoniter.API,
	//%: {{- end }}
	//%: {{- if .Yafuds }}
	yafudsClient yafuds.Client,
	validator *validation.Validator,
	//%: {{- end }}
	//%: {{- if and .Geb .Examples }}
	gebQueue *geb.Queue,
	//%: {{- end }}
	//%: {{- if .SiteConfig }}
	siteConfigClient *client.SiteConfig,
	//%: {{- end }}
) *Service {
	return &Service{
		//%: {{- if .Centrifuge }}
		centrifugeClient: centrifugeClient,
		centrifugeJSON:   centrifugeJSON,
		//%: {{- end }}
		//%: {{- if .Yafuds }}
		yafudsClient: yafudsClient,
		validator:    validator,
		//%: {{- end }}
		//%: {{- if and .Geb .Examples }}
		gebQueue: gebQueue,
		//%: {{- end }}
		//%: {{- if .SiteConfig }}
		siteConfigClient: siteConfigClient,
		//%: {{- end }}
	}
}

//%: {{ if .Examples }}{{ `
func (s *Service) Dummy(ctx context.Context, req *skeleton.DummyRequest) error {
	//%: ` | replace "skeleton" .SchemaPackage | trim }}
	//%: {{- if .Centrifuge }}{{ `
	s.SendCentrifuge(ctx, centrifuge.GlobalChannel("dummy", "dummy_site_group", "dummy"), skeleton.CentrifugeData{
		Group:     req.DummyGroup,
		DummyUUID: req.DummyUUID,
	})
	//%: ` | replace "skeleton" .SchemaPackage }}{{- end }}

	return nil
} //%: {{ end }}

//%: {{- if and .Centrifuge .Examples }}
// todo: remove
//  Centrifuge example
func (s *Service) SendCentrifuge(ctx context.Context, centrifugeChannel string, eventData interface{}) {
	data, err := s.centrifugeJSON.Marshal(eventData)
	if err != nil {
		//%:{{ `
		err = skeleton.SemanticError{Err: err, Msg: fmt.Sprintf("unable to marshal eventData of type: %T", eventData)}.E()
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
		log.Error(ctx, err.Error(), "error", err)
		return
	}

	err = s.centrifugeClient.Publish(ctx, &centrifuge.PublishRequest{
		Channel: centrifugeChannel,
		Data:    data,
	})
	if err != nil {
		log.Error(ctx, err.Error(), "error", err)
	}
} //%: {{- end }}

//%: {{- if and .Yafuds .Examples }}
// todo: remove
//  Update example
//%:{{ `
func (s *Service) Update(ctx context.Context, req *skeleton.UpdateDummyRequest, keys map[string]interface{}) (*skeleton.DummyType, error) {
	//%: ` | replace "skeleton" .SchemaPackage | trim }}
	var (
		//%:{{ `
		dum *skeleton.DummyType
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
		changedFields []string
		oldVersion    string
		newVersion    string
	)

	var err = s.yafudsClient.Retry(ctx, func() error {
		//%:{{ `
		dum = &skeleton.DummyType{}
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
		columnsData := dummyToColumns(dum)

		cells, err := s.yafudsClient.GetLatestCells(ctx, req.DummyUUID.String(), columnsData)
		if err != nil {
			return yafudsError{Err: err}.E()
		}
		dum.DummyUUID = req.DummyUUID

		changedFields = make([]string, 0, len(keys))
		changedColumns := make(map[string]interface{}, len(columnsData))

		merges := []struct {
			col string
			dst interface{}
			src interface{}
		}{
			//%:{{ `
			{col: skeleton.YcolBaseEditable, dst: &dum.BaseEditable, src: &req.BaseEditable},
			//%: ` | replace "skeleton" .SchemaPackage | trim }}
		}

		for _, m := range merges {
			cf, err := shallow.Merge(m.dst, m.src, keys)
			if err != nil {
				//%:{{ `
				return skeleton.SemanticError{Err: err}.E()
				//%: ` | replace "skeleton" .SchemaPackage | trim }}
			}
			if len(cf) > 0 {
				changedColumns[m.col] = columnsData[m.col]
				changedFields = append(changedFields, cf...)
			}
		}

		if len(changedColumns) == 0 && len(changedFields) > 0 {
			//%:{{ `
			return skeleton.SemanticError{Msg: "changedColumns is empty but changeFields is not"}.E()
			//%: ` | replace "skeleton" .SchemaPackage | trim }}
		}

		if len(changedFields) == 0 && len(changedColumns) > 0 {
			//%:{{ `
			return skeleton.SemanticError{Msg: "changeFields is empty but changedColumns is not"}.E()
			//%: ` | replace "skeleton" .SchemaPackage | trim }}
		}

		if len(changedFields) == 0 {
			return nil
		}

		if err := s.validator.Validate(dum); err != nil {
			return err
		}

		colRefs, err := s.yafudsClient.PutCells(ctx, req.DummyUUID.String(), cells, changedColumns)
		if err != nil {
			return err
		}

		oldVersion, newVersion = versionOldNew(cells, colRefs)

		return nil
	})
	if err != nil {
		//%:{{ `
		if schema.ErrorCode(err) == skeleton.ErrValidation {
			//%: ` | replace "skeleton" .SchemaPackage | trim }}
			return nil, err
		}

		//%:{{ `
		return nil, skeleton.SemanticError{Err: err}.E()
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
	}

	//%: {{ if .Geb }}
	if len(changedFields) > 0 {
		s.publishUpdated(ctx, dum.DummyUUID, changedFields, oldVersion, newVersion)
	}
	//%: {{- end }}

	return dum, nil
}

// todo: remove
func versionOldNew(oldCells []*yafuds.Cell, newColRefs map[string]uint32) (string, string) {
	colRefs := make(map[string]uint32, len(oldCells)+len(newColRefs))
	for _, c := range oldCells {
		colRefs[c.ColumnKey()] = c.RefKey()
	}

	oldPairs := make([]string, 0, len(colRefs))
	for col, ref := range colRefs {
		oldPairs = append(oldPairs, fmt.Sprintf("%s:%d", col, ref))
	}
	sort.Strings(oldPairs)

	for col, ref := range newColRefs {
		colRefs[col] = ref
	}

	newPairs := make([]string, 0, len(colRefs))
	for col, ref := range colRefs {
		newPairs = append(newPairs, fmt.Sprintf("%s:%d", col, ref))
	}
	sort.Strings(newPairs)

	return strings.Join(oldPairs, ","), strings.Join(newPairs, ",")
}

//%:{{ `
func dummyToColumns(dum *skeleton.DummyType) map[string]interface{} {
	//%: ` | replace "skeleton" .SchemaPackage | trim }}
	return map[string]interface{}{
		//%:{{ `
		skeleton.YcolBaseEditable: &dum.BaseEditable,
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
	}
} //%: {{- end }}

//%: {{- if and .Geb .Examples }}
// todo: remove
//  Geb example
func (s *Service) publishUpdated(
	ctx context.Context,
	dummyUUID uuid.UUID,
	changedFields []string,
	oldVersion string,
	newVersion string,
) {
	s.publishEvent(
		ctx,
		//%:{{ `
		"event/skeleton-service/dummy/updated/v1",
		//%: ` | replace "skeleton-service" .ProjectName | trim }}{{ `
		&skeleton.UpdatedEvent{
			//%: ` | replace "skeleton" .SchemaPackage | trim }}
			UUID:          dummyUUID,
			ChangedFields: changedFields,
			OldVersion:    oldVersion,
			NewVersion:    newVersion,
		},
	)
}

// todo: remove
//  Geb example
func (s *Service) publishEvent(ctx context.Context, eventName string, eventData interface{}) {
	err := s.gebQueue.Publish(eventName).
		Context(ctx).
		Body(eventData).
		Do()

	if err != nil {
		err = gebError{Err: err}.E()
		log.Warn(ctx, err.Error(), "error", err)
	}
}

//%: {{- end }}
