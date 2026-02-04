package fields

import (
	"encoding/json"
)

// RelationshipSerialization handles JSON serialization for relationships
type RelationshipSerialization interface {
	// SerializeRelationship serializes a relationship to JSON
	SerializeRelationship(item interface{}) (map[string]interface{}, error)

	// SerializeRelationships serializes multiple relationships to JSON
	SerializeRelationships(items []interface{}) ([]map[string]interface{}, error)
}

// RelationshipSerializationImpl implements RelationshipSerialization
type RelationshipSerializationImpl struct {
	field RelationshipField
}

// NewRelationshipSerialization creates a new relationship serialization handler
func NewRelationshipSerialization(field RelationshipField) *RelationshipSerializationImpl {
	return &RelationshipSerializationImpl{
		field: field,
	}
}

// SerializeRelationship serializes a relationship to JSON
func (rs *RelationshipSerializationImpl) SerializeRelationship(item interface{}) (map[string]interface{}, error) {
	if item == nil {
		return map[string]interface{}{
			"value": nil,
		}, nil
	}

	// Convert item to JSON-compatible format
	jsonData := map[string]interface{}{
		"type":     rs.field.GetRelationshipType(),
		"name":     rs.field.GetRelationshipName(),
		"resource": rs.field.GetRelatedResource(),
		"value":    item,
	}

	return jsonData, nil
}

// SerializeRelationships serializes multiple relationships to JSON
func (rs *RelationshipSerializationImpl) SerializeRelationships(items []interface{}) ([]map[string]interface{}, error) {
	if items == nil || len(items) == 0 {
		return []map[string]interface{}{}, nil
	}

	serialized := make([]map[string]interface{}, 0, len(items))

	for _, item := range items {
		jsonData, err := rs.SerializeRelationship(item)
		if err != nil {
			return nil, err
		}
		serialized = append(serialized, jsonData)
	}

	return serialized, nil
}

// ToJSON converts relationship to JSON string
func (rs *RelationshipSerializationImpl) ToJSON(item interface{}) (string, error) {
	jsonData, err := rs.SerializeRelationship(item)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// ToJSONArray converts relationships to JSON array string
func (rs *RelationshipSerializationImpl) ToJSONArray(items []interface{}) (string, error) {
	jsonData, err := rs.SerializeRelationships(items)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
