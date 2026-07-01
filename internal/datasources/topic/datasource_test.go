package topic_test

import (
	"context"
	"testing"

	datasourcetopic "github.com/gravitino/terraform-provider-gravitino/internal/datasources/topic"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestTopicDataSourceMetadata(t *testing.T) {
	d := datasourcetopic.NewTopicDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_topic" {
		t.Errorf("Expected type name gravitino_topic, got %s", resp.TypeName)
	}
}

func TestTopicsDataSourceMetadata(t *testing.T) {
	d := datasourcetopic.NewTopicsDataSource()
	var req datasource.MetadataRequest
	var resp datasource.MetadataResponse
	d.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "gravitino_topics" {
		t.Errorf("Expected type name gravitino_topics, got %s", resp.TypeName)
	}
}
