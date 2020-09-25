package storage

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"gitlab.com/proemergotech/yafuds-client-go/client"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
)

const (
	YcolBaseEditable = "BASE_EDITABLE"
)

type Yafuds struct {
	client client.Client
}

func NewYafuds(client client.Client) *Yafuds {
	return &Yafuds{
		client: client,
	}
}

// todo: remove
//   Implementation example for save value
func (y *Yafuds) Save(ctx context.Context, dummy *service.DummyType) (string, string, error) {
	var oldVersion, newVersion string
	err := y.client.Retry(ctx, func() error {
		columns := []client.Column{
			{ColumnKey: YcolBaseEditable},
		}

		cells, err := y.client.GetCells(ctx, dummy.UUID.String(), columns)
		if err != nil {
			return err
		}

		bodies := map[string]interface{}{
			YcolBaseEditable: dummy.BaseEditable,
		}

		colRefs, err := y.client.PutCells(ctx, dummy.UUID.String(), cells, bodies)
		if err != nil {
			return err
		}

		oldVersion, newVersion = versionOldNew(cells, colRefs)

		return nil
	})

	if err != nil {
		return "", "", yafudsError{Err: err}.E()
	}

	return oldVersion, newVersion, nil
}

func versionOldNew(oldCells []*client.Cell, newColRefs map[string]uint32) (string, string) {
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
