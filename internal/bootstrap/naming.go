package bootstrap

import coredoctype "github.com/galaxylabs/gogal-studio/internal/core/doctype"

func DocTypeTableName(doctypeName string) string {
	tableName, err := coredoctype.TableName(doctypeName)
	if err != nil {
		return ""
	}

	return tableName
}
