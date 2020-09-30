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
	Name() string
	SetName(name string)
	SetSelf(self Strategy) error
	//Setup(mode TradeMode, exchanges ...Exchange) error
	Setup(mode TradeMode, exchanges ...interface{}) error
	IsStopped() bool
	StopNow()
	TradeMode() TradeMode
	SetOptions(options map[string]interface{}) error
	Run() error
	OnInit() error
	OnTick() error
	OnExit() error
}

// StrategyBase Strategy base class
type StrategyBase struct {
	self      interface{}
	name      string
	tradeMode TradeMode
	Exchanges []Exchange
	Exchange  Exchange
	stopped   bool
}

// SetSelf 设置 self 对象
func (s *StrategyBase) SetSelf(self Strategy) error {
	s.self = self.(interface{})
	return nil
}

// Setup Setups the exchanges
func (s *StrategyBase) Setup(mode TradeMode, exchanges ...interface{}) error { // Exchange
	if len(exchanges) == 0 {
		return fmt.Errorf("no exchanges")
	}
	s.tradeMode = mode
	for _, v := range exchanges {
		ex, ok := v.(Exchange)
		if !ok {
			return fmt.Errorf("Exchange Only")
		}
		s.Exchanges = append(s.Exchanges, ex)
	}
	s.Exchange = s.Exchanges[0]
	s.stopped = false
	return nil
}

// SetOptions Sets the options for the strategy
func (s *StrategyBase) SetOptions(options map[string]interface{}) error {
	return setOptions(s.self, options)
}

// GetOptions Returns the options of strategy
func (s *StrategyBase) GetOptions() (optionMap map[string]*StrategyOption) {
	return getOptions(s.self)
}

func (s *StrategyBase) TradeMode() TradeMode {
	return s.tradeMode
}

func (s *StrategyBase) IsStopped() bool {
	return s.stopped
}

func (s *StrategyBase) StopNow() {
	s.stopped = true
}

func (s *StrategyBase) SetName(name string) {
	s.name = name
}

func (s *StrategyBase) Name() string {
	return s.name
}

// SpotStrategyBase Strategy base class
type SpotStrategyBase struct {
	self      interface{}
	name      string
	tradeMode TradeMode
	Exchanges []SpotExchange
	Exchange  SpotExchange
}

// SetSelf 设置 self 对象
func (s *SpotStrategyBase) SetSelf(self Strategy) error {
	s.self = self.(interface{})
	return nil
}

// Setup Setups the exchanges
func (s *SpotStrategyBase) Setup(mode TradeMode, exchanges ...interface{}) error { // Exchange
	if len(exchanges) == 0 {
		return fmt.Errorf("no exchanges")
	}
	s.tradeMode = mode
	for _, v := range exchanges {
		ex, ok := v.(SpotExchange)
		if !ok {
			return fmt.Errorf("SpotExchange only")
		}
		s.Exchanges = append(s.Exchanges, ex)
	}
	s.Exchange = s.Exchanges[0]
	return nil
}

// SetOptions Sets the options for the strategy
func (s *SpotStrategyBase) SetOptions(options map[string]interface{}) error {
	return setOptions(s.self, options)
}

// GetOptions Returns the options of strategy
func (s *SpotStrategyBase) GetOptions() (optionMap map[string]*StrategyOption) {
	return getOptions(s.self)
}

func (s *SpotStrategyBase) TradeMode() TradeMode {
	return s.tradeMode
}

func (s *SpotStrategyBase) SetName(name string) {
	s.name = name
}

func (s *SpotStrategyBase) Name() string {
	return s.name
}

// 组合策略，期现等
// CStrategyBase Strategy base class
type CStrategyBase struct {
	self          interface{}
	name          string
	tradeMode     TradeMode
	Exchanges     []Exchange
	SpotExchanges []SpotExchange
	stopped       bool
}

// SetSelf 设置 self 对象
func (s *CStrategyBase) SetSelf(self Strategy) error {
	s.self = self.(interface{})
	return nil
}

// Setup Setups the exchanges
func (s *CStrategyBase) Setup(mode TradeMode, exchanges ...interface{}) error { // Exchange
	if len(exchanges) == 0 {
		return fmt.Errorf("no exchanges")
	}
	s.tradeMode = mode
	for _, v := range exchanges {
		if ex, ok := v.(Exchange); ok {
			s.Exchanges = append(s.Exchanges, ex)
			continue
		}

		if ex, ok := v.(SpotExchange); ok {
			s.SpotExchanges = append(s.SpotExchanges, ex)
		}
	}
	s.stopped = false
	return nil
}

// SetOptions Sets the options for the strategy
func (s *CStrategyBase) SetOptions(options map[string]interface{}) error {
	return setOptions(s.self, options)
}

// GetOptions Returns the options of strategy
func (s *CStrategyBase) GetOptions() (optionMap map[string]*StrategyOption) {
	return getOptions(s.self)
}

func (s *CStrategyBase) TradeMode() TradeMode {
	return s.tradeMode
}

func (s *CStrategyBase) IsStopped() bool {
	return s.stopped
}

func (s *CStrategyBase) StopNow() {
	s.stopped = true
}

func (s *CStrategyBase) SetName(name string) {
	s.name = name
}

func (s *CStrategyBase) Name() string {
	return s.name
}

// SetOptions Sets the options for the strategy
func setOptions(s interface{}, options map[string]interface{}) error {
	if len(options) == 0 {
		return nil
	}

	rawOptionsOrigin := getOptions(s)
	rawOptions := map[string]*StrategyOption{}

	for k, v := range rawOptionsOrigin {
		key := strings.ReplaceAll(strings.ToLower(k), "_", "")
		rawOptions[key] = v
	}

	// 反射成员变量
	val := reflect.ValueOf(s)

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

func getOptions(s interface{}) (optionMap map[string]*StrategyOption) {
	//log.Info("GetOptions")
	optionMap = map[string]*StrategyOption{}

	if s == nil {
		return
	}

	val := reflect.ValueOf(s)

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
		defaultValue := getDefaultValue(fieldKind, defaultValueString)

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

func getDefaultValue(kind reflect.Kind, value string) interface{} {
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
