package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testObj struct {
	Time ADCTime   `json:"time"`
	ID   ADCInt64  `json:"id"`
	Name ADCString `json:"name"`
}

func TestJSONSerDes(t *testing.T) {
	tt := NewADCTime("2019-02-27 18:49:15")
	assert.NotNil(t, tt)

	id := NewADCInt64(100)
	name := NewADCString("192.168.100.10")
	obj := testObj{Time: tt, ID: id, Name: name}

	res, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
	}

	var objBack testObj
	err = json.Unmarshal(res, &objBack)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, objBack.Time, obj.Time)
	assert.Equal(t, objBack.ID, obj.ID)
	assert.Equal(t, objBack.Name, obj.Name)
}
