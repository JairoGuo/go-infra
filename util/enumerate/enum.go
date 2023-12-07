package enumerate

type Enum[T any] interface {
	Value() T
	String() string
}

type Int interface {
	~int
}

type EnumContainer[K Int, T any] struct {
	dataMap map[K]T
}

func (container EnumContainer[K, T]) Value(enumKey K, value T) {
	container.dataMap[enumKey] = value
}

func (container EnumContainer[K, T]) Get(enumKey K) T {
	return container.dataMap[enumKey]
}

func CreateEnum[K Int, T any]() EnumContainer[K, T] {
	container := EnumContainer[K, T]{}
	return container
}
