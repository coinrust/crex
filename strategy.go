package crex

import (
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"strings"
)

const (
	// StrategyOptionTag 选项Tag
	StrategyOptionTag = "opt"
)

// StrategyOption 策略参数
type StrategyOption struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Type         string      `json:"type"`
	Value        interface{} `json:"value"`
	DefaultValue interface{} `json:"default_value"`
}

// Strategy interface
type Strategy interface {
	SetSelf(self Strategy) error
	Setup(mode TradeMode, exchanges ...Exchange) error
	TradeMode() TradeMode
	SetOptions(options map[string]interface{}) error
	Run() error
	OnInit() error
	OnTick() error
	OnDeinit() error
}

// StrategyBase Strategy base class
type StrategyBase struct {
	self      interface{}
	tradeMode TradeMode
	Exchanges []Exchange
	Exchange  Exchange
}

// SetSelf 设置 self 对象
func (s *StrategyBase) SetSelf(self Strategy) error {
	s.self = self.(interface{})
	return nil
}

// Setup Setups the exchanges
func (s *StrategyBase) Setup(mode TradeMode, exchanges ...Exchange) error {
	if len(exchanges) == 0 {
		return fmt.Errorf("no exchanges")
	}
	s.tradeMode = mode
	s.Exchanges = append(s.Exchanges, exchanges...)
	s.Exchange = exchanges[0]
	return nil
}

// SetOptions Sets the options for the strategy
func (s *StrategyBase) SetOptions(options map[string]interface{}) error {
	if len(options) == 0 {
		return nil
	}

	rawOptionsOrigin := s.GetOptions()
	rawOptions := map[string]*StrategyOption{}

	for k, v := range rawOptionsOrigin {
		key := strings.ReplaceAll(strings.ToLower(k), "_", "")
		rawOptions[key] = v
	}

	// 反射成员变量
	val := reflect.ValueOf(s.self)

	// If it's an interface or a pointer, unwrap it.
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		return nil
	}

	for name, value := range options {
		var fieldName string

		key := strings.ReplaceAll(name, "_", "")
		if ipi, ok := rawOptions[key]; !ok {
			continue
		} else {
			fieldName = ipi.Name
		}

		//fmt.Println(fieldName)

		v := val.FieldByName(fieldName)
		if !v.IsValid() {
			continue
		}

		switch v.Kind() {
		default:
			fmt.Printf("Error Kind: %v\n", v.Kind())
		case reflect.Bool:
			v.SetBool(cast.ToBool(value))
		case reflect.String:
			v.SetString(value.(string))
		case reflect.Int:
			v.SetInt(cast.ToInt64(value))
		case reflect.Int8:
			v.SetInt(cast.ToInt64(value))
		case reflect.Int16:
			v.SetInt(cast.ToInt64(value))
		case reflect.Int32:
			v.SetInt(cast.ToInt64(value))
		case reflect.Int64:
			v.SetInt(cast.ToInt64(value))
		case reflect.Uint:
			v.SetUint(cast.ToUint64(value))
		case reflect.Uint8:
			v.SetUint(cast.ToUint64(value))
		case reflect.Uint16:
			v.SetUint(cast.ToUint64(value))
		case reflect.Uint32:
			v.SetUint(cast.ToUint64(value))
		case reflect.Uint64:
			v.SetUint(cast.ToUint64(value))
		case reflect.Float32:
			v.SetFloat(cast.ToFloat64(value))
		case reflect.Float64:
			v.SetFloat(cast.ToFloat64(value))
			// case reflect.Struct:
			// 	v.Set(reflect.ValueOf(value))
		}
	}

	return nil
}

// GetOptions Returns the options of strategy
func (s *StrategyBase) GetOptions() (optionMap map[string]*StrategyOption) {
	//log.Info("GetOptions")
	optionMap = map[string]*StrategyOption{}

	if s.self == nil {
		return
	}

	val := reflect.ValueOf(s.self)

	// If it's an interface or a pointer, unwrap it.
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		return
	}

	valNumFields := val.NumField()

	for i := 0; i < valNumFields; i++ {
		field := val.Field(i)
		fieldKind := field.Kind()
		if fieldKind == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			continue
		}

		typeField := val.Type().Field(i)
		fieldName := typeField.Name
		tag := typeField.Tag

		if !field.CanInterface() {
			continue
		}

		option := tag.Get(StrategyOptionTag)

		if option == "" {
			continue
		}

		var description string
		var defaultValueString string
		index := strings.Index(option, ",")
		//fmt.Printf("tag: %v i: %v\n", option, index)
		if index != -1 {
			description = option[0:index]
			defaultValueString = option[index+1:]
		} else {
			description = option
		}
		value := field.Interface()
		defaultValue := s.getDefaultValue(fieldKind, defaultValueString)

		optionMap[fieldName] = &StrategyOption{
			Name:         fieldName,
			Description:  description,
			Type:         typeField.Type.String(),
			Value:        value,
			DefaultValue: defaultValue,
		}
		//log.Infof("F: %v V: %v", fieldName, value)
	}

	return
}

func (s *StrategyBase) getDefaultValue(kind reflect.Kind, value string) interface{} {
	switch kind {
	case reflect.Bool:
		return cast.ToBool(value)
	case reflect.String:
		return value
	case reflect.Int:
		return cast.ToInt(value)
	case reflect.Int8:
		return cast.ToInt8(value)
	case reflect.Int16:
		return cast.ToInt16(value)
	case reflect.Int32:
		return cast.ToInt32(value)
	case reflect.Int64:
		return cast.ToInt64(value)
	case reflect.Uint:
		return cast.ToUint(value)
	case reflect.Uint8:
		return cast.ToUint8(value)
	case reflect.Uint16:
		return cast.ToUint16(value)
	case reflect.Uint32:
		return cast.ToUint32(value)
	case reflect.Uint64:
		return cast.ToInt64(value)
	case reflect.Float32:
		return cast.ToFloat32(value)
	case reflect.Float64:
		return cast.ToFloat64(value)
	}
	return 0
}

func (s *StrategyBase) TradeMode() TradeMode {
	return s.tradeMode
}
