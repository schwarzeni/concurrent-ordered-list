package concurrentorderdlist

type List interface {
	// Contains 检查一个元素是否存在，如果存在则返回 true，否则返回 false
	Contains(value int) bool

	// Insert 插入一个元素，如果此操作成功插入一个元素，则返回 true，否则返回 false
	Insert(value int) bool

	// Delete 删除一个元素，如果此操作成功删除一个元素，则返回 true，否则返回 false
	Delete(value int) bool

	// Range 遍历此有序链表的所有元素，如果 f 返回 false，则停止遍历
	Range(f func(value int) bool)

	// Len 返回有序链表的元素个数
	Len() int
}
