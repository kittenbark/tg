package tg

import (
	"encoding/json"
	"testing"
)

func TestRichBlockBlockQuotationNestedBlocks(t *testing.T) {
	data := []byte(`{
		"type": "blockquote",
		"blocks": [
			{
				"type": "blockquote",
				"blocks": [
					{"type": "paragraph", "text": "hi"}
				]
			}
		]
	}`)

	var block RichBlockBlockQuotation
	if err := json.Unmarshal(data, &block); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if len(block.Blocks) != 1 {
		t.Fatalf("expected 1 outer block, got %d", len(block.Blocks))
	}
	inner, ok := block.Blocks[0].(*RichBlockBlockQuotation)
	if !ok {
		t.Fatalf("expected *RichBlockBlockQuotation, got %T", block.Blocks[0])
	}
	if len(inner.Blocks) != 1 {
		t.Fatalf("expected 1 inner block, got %d", len(inner.Blocks))
	}
	if _, ok := inner.Blocks[0].(*RichBlockParagraph); !ok {
		t.Fatalf("expected *RichBlockParagraph, got %T", inner.Blocks[0])
	}
}
