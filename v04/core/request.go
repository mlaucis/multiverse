package core

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/errmsg"
)

type (
	// RequestCondition holds down the possible filtering values for the fields
	RequestCondition struct {
		Eq  interface{}   `json:"eq,omitempty"`
		Neq interface{}   `json:"neq,omitempty"`
		Lt  interface{}   `json:"lt,omitempty"`
		Lte interface{}   `json:"lte,omitempty"`
		Gt  interface{}   `json:"gt,omitempty"`
		Gte interface{}   `json:"gte,omitempty"`
		In  []interface{} `json:"in,omitempty"`
		Nin []interface{} `json:"nin,omitempty"`
	}

	// EventCondition holds the possible event fields to be filtered
	EventCondition struct {
		Language *RequestCondition             `json:"language,omitempty"`
		Location *RequestCondition             `json:"location,omitempty"`
		Metadata *map[string]*RequestCondition `json:"metadata,omitempty"`
		ObjectID *RequestCondition             `json:"tg_object_id"`
		Owned    *RequestCondition
		Priority *RequestCondition `json:"priority,omitempty"`
		Type     *RequestCondition `json:"type,omitempty"`
		Object   *struct {
			ID   *RequestCondition `json:"id,omitempty"`
			Type *RequestCondition `json:"type,omitempty"`
		} `json:"object,omitempty"`
		Target *struct {
			ID   *RequestCondition `json:"id,omitempty"`
			Type *RequestCondition `json:"type,omitempty"`
		} `json:"target,omitempty"`
	}
)

// condition generate the condition out of a field
func (s *RequestCondition) condition(fieldName string, paramID int) (cond string, params []interface{}, paramIDStop int, err error) {
	condition := []string{}

	getFieldType := func(fieldValue interface{}) (string, error) {
		fieldType := ""
		switch fieldValue.(type) {
		case bool:
			fieldType = "BOOLEAN"
		case int:
			fieldType = "BIGINT"
		case int64:
			fieldType = "BIGINT"
		case float64:
			// Unfortunately here is where the JSON spec or parser go against common sense and default to FLOAT
			// so we need to force this upon our type to make sure we really have a float there
			// I'm really sorry CPU for what I'm doing to you here
			if math.Trunc(fieldValue.(float64)) == fieldValue.(float64) {
				fieldType = "BIGINT"
			} else {
				fieldType = "FLOAT"
			}
		case string:
			fieldType = "TEXT"
		default:
			return "", errmsg.ErrInvalidFieldTypeError
		}

		return fieldType, nil
	}

	getFieldCondition := func(fieldName, operation string, fieldValue interface{}, paramID int) (string, error) {
		fieldType, err := getFieldType(fieldValue)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf(`(%s)::%s %s $%d::%s`, fieldName, fieldType, operation, paramID, fieldType), nil
	}

	if s.Eq != interface{}(nil) {
		cond, er := getFieldCondition(fieldName, "=", s.Eq, paramID)
		if er == nil {
			condition = append(condition, cond)
			params = append(params, s.Eq)
			paramID++
		}
	}

	if s.Neq != interface{}(nil) {
		cond, er := getFieldCondition(fieldName, "<>", s.Neq, paramID)
		if er == nil {
			condition = append(condition, cond)
			params = append(params, s.Neq)
			paramID++
		}
	}

	if s.Lt != interface{}(nil) {
		cond, er := getFieldCondition(fieldName, "<", s.Lt, paramID)
		if er == nil {
			condition = append(condition, cond)
			params = append(params, s.Lt)
			paramID++
		}
	}

	if s.Lte != interface{}(nil) {
		cond, er := getFieldCondition(fieldName, "<=", s.Lte, paramID)
		if er == nil {
			condition = append(condition, cond)
			params = append(params, s.Lte)
			paramID++
		}
	}

	if s.Gt != interface{}(nil) {
		cond, er := getFieldCondition(fieldName, ">", s.Gt, paramID)
		if er == nil {
			condition = append(condition, cond)
			params = append(params, s.Gt)
			paramID++
		}
	}

	if s.Gte != interface{}(nil) {
		cond, er := getFieldCondition(fieldName, ">=", s.Gte, paramID)
		if er == nil {
			condition = append(condition, cond)
			params = append(params, s.Gte)
			paramID++
		}
	}

	if len(s.In) > 0 {
		inType, err := getFieldType(s.In[0])
		if err == nil {
			buffer := []string{fmt.Sprintf("$%d::%s", paramID, inType)}
			paramID++
			for i := 1; i < len(s.In); i++ {
				if fieldType, err := getFieldType(s.In[i]); err == nil && fieldType == inType {
					buffer = append(buffer, fmt.Sprintf("$%d::%s", paramID, inType))
					paramID++
				}
			}

			condition = append(condition, fmt.Sprintf("(%s)::%s IN (%s)", fieldName, inType, strings.Join(buffer, ", ")))
			params = append(params, s.In)
		}
	}

	if len(s.Nin) > 0 {
		ninType, err := getFieldType(s.Nin[0])
		if err == nil {
			buffer := []string{fmt.Sprintf("$%d::%s", paramID, ninType)}
			paramID++
			for i := 1; i < len(s.In); i++ {
				if fieldType, err := getFieldType(s.Nin[i]); err != nil && fieldType == ninType {
					buffer = append(buffer, fmt.Sprintf("$%d::%s", paramID, ninType))
					paramID++
				}
			}

			condition = append(condition, fmt.Sprintf("(%s)::%s IN (%s)", fieldName, ninType, strings.Join(buffer, ", ")))
			params = append(params, s.Nin)
		}
	}

	if len(condition) == 0 {
		return "", []interface{}{}, 0, errmsg.ErrInvalidConditionLengthError
	}

	cond = strings.Join(condition, " AND ")

	return cond, params, paramID, nil
}

// conditions generates the conditions for events
func (e *EventCondition) conditions(startPlaceholderID int) (query string, params []interface{}, paramIDStop int, err error) {
	qry := []string{}

	paramID := startPlaceholderID

	checkSimpleField := func(field *RequestCondition, fieldStr string) error {
		if field == nil {
			return nil
		}
		cond, param, paramIDNext, er := field.condition(fieldStr, paramID)
		if er != nil {
			return er
		}
		if cond != "" {
			qry = append(qry, cond)
			params = append(params, param...)
			paramID = paramIDNext
		}

		return nil
	}

	if err := checkSimpleField(e.Type, `json_data->>'type'`); err != nil {
		return "", []interface{}{}, 0, err
	}

	if err := checkSimpleField(e.Language, `json_data->>'language'`); err != nil {
		return "", []interface{}{}, 0, err
	}

	if err := checkSimpleField(e.Priority, `json_data->>'priority'`); err != nil {
		return "", []interface{}{}, 0, err
	}

	if err := checkSimpleField(e.Location, `json_data->>'location'`); err != nil {
		return "", []interface{}{}, 0, err
	}

	if err := checkSimpleField(e.ObjectID, `json_data->>'object_id'`); err != nil {
		return "", []interface{}{}, 0, err
	}

	if e.Object != nil {
		if err := checkSimpleField(e.Object.ID, `json_data->'object'->>'id'`); err != nil {
			return "", []interface{}{}, 0, err
		}
		if err := checkSimpleField(e.Object.Type, `json_data->'object'->>'type'`); err != nil {
			return "", []interface{}{}, 0, err
		}
	}

	if e.Target != nil {
		if err := checkSimpleField(e.Target.ID, `json_data->'target'->>'id'`); err != nil {
			return "", []interface{}{}, 0, err
		}
		if err := checkSimpleField(e.Target.Type, `json_data->'target'->>'type'`); err != nil {
			return "", []interface{}{}, 0, err
		}
	}

	if e.Metadata != nil {
		for idx, val := range *e.Metadata {
			if err := checkSimpleField(val, fmt.Sprintf(`json_data->'metadata'->>'%s'`, idx)); err != nil {
				return "", []interface{}{}, 0, err
			}
		}
	}

	return strings.Join(qry, " AND "), params, paramID, nil
}

// Process will process the current event filter and retrieve the condition and the parameters for the condition
func (e *EventCondition) Process(startingPlaceholderID int) (requestCondition string, requestParams []interface{}, err []errors.Error) {
	if e == nil {
		return "", []interface{}{}, nil
	}

	var er error
	requestCondition, requestParams, _, er = e.conditions(startingPlaceholderID)
	if err != nil {
		return "", nil, []errors.Error{errmsg.ErrInvalidFilterConditionError.UpdateInternalMessage(er.Error())}
	}

	if requestCondition != "" {
		requestCondition = "AND " + requestCondition
	}

	return requestCondition, requestParams, nil
}

// NewEventFilter will create a new event filter out of a filter string
func NewEventFilter(filter string) (*EventCondition, []errors.Error) {
	result := &EventCondition{}
	if filter == "" || filter == "{}" {
		return nil, nil
	}

	err := json.Unmarshal([]byte(filter), result)
	if err != nil {
		return nil, []errors.Error{errmsg.ErrInvalidFilterConditionError.UpdateInternalMessage(err.Error())}
	}

	return result, nil
}
