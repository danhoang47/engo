package fiber

import (
	"engo/internal/protocol"
	"engo/pkg/engine"
	"engo/pkg/scene"
	"time"
)

type Reconciler struct {
	CurrentRoot    *Fiber // Cây đang hiển thị
	WipRoot        *Fiber // Cây đang xây dựng (Work In Progress)
	NextUnitOfWork *Fiber // Con trỏ "Cursor" hiện tại

	engine interface{}
}

func NewReconciler(engine interface{}) *Reconciler {
	return &Reconciler{
		CurrentRoot:    nil,
		WipRoot:        nil,
		NextUnitOfWork: nil,

		engine: engine,
	}
}

// ScheduleUpdate use to mark a node is dirty
// and need repaint action
//
// rootNode is page-level node of scene-graph
func (r *Reconciler) ScheduleUpdate(rootNode *scene.Node) {
	if r.WipRoot != nil {
		r.NextUnitOfWork = nil
	}

	r.WipRoot = CreateWorkInProgress(r.CurrentRoot, rootNode.Props)
	r.WipRoot.Node = rootNode
	r.NextUnitOfWork = r.WipRoot
}

func (r *Reconciler) WorkLoop(budgetMs int64) bool {
	startTime := time.Now().UnixMilli()

	for r.NextUnitOfWork != nil {
		if budgetMs > 0 {
			currentTime := time.Now().UnixMilli()
			if (currentTime - startTime) >= budgetMs {
				return true
			}
		}

		r.NextUnitOfWork = r.performUnitOfWork(r.NextUnitOfWork)
	}

	if r.WipRoot != nil {
		r.CommitRoot()
	}

	return false
}

func (r *Reconciler) performUnitOfWork(unit *Fiber) *Fiber {
	next := r.beginWork(unit)

	if next != nil {
		return next
	}

	curr := unit
	for curr != nil {
		r.completeWork(curr)

		if curr.Sibling != nil {
			return curr.Sibling // Sang ngang
		}

		curr = curr.Parent
	}

	return nil
}

func (r *Reconciler) beginWork(fiber *Fiber) *Fiber {
	node := fiber.Node // Source of Truth

	parentMatrix := protocol.IdentityMatrix()
	if fiber.Parent != nil {
		parentMatrix = fiber.Parent.GlobalMatrix
	}

	fiber.ComputeMatrix(parentMatrix)

	if fiber.Alternate != nil &&
		(node.Flags&scene.FlagLayoutDirty == 0) &&
		(node.Flags&scene.FlagSubtreeDirty == 0) {
		return r.bailoutOnAlreadyFinishedWork(fiber, fiber.Alternate)
	}

	var sceneChildren []*scene.Node
	if node != nil {
		sceneChildren = node.Children
	}

	r.reconcileChildren(fiber, sceneChildren)

	return fiber.Child
}

func (r *Reconciler) bailoutOnAlreadyFinishedWork(current *Fiber, alternate *Fiber) *Fiber {
	child := alternate.Child

	if child != nil {
		current.Child = child
		child.Parent = current
	}

	return child
}

func (r *Reconciler) reconcileChildren(returnFiber *Fiber, newChildren []*scene.Node) {
	existingChildren := make(map[uint32]*Fiber)

	// first child of current wip fiber
	var currentOldFiber *Fiber
	if returnFiber.Alternate != nil {
		currentOldFiber = returnFiber.Alternate.Child
	}

	temp := currentOldFiber
	for temp != nil {
		if temp.Key != 0 {
			existingChildren[temp.Key] = temp
		} else {
			existingChildren[uint32(temp.Node.ID)] = temp
		}
		temp = temp.Sibling
	}

	// 2. Duyệt qua danh sách con MỚI (Scene Graph)
	var prevSibling *Fiber = nil
	var resultingFirstChild *Fiber = nil

	for i, newNode := range newChildren {
		matchedFiber, exists := existingChildren[newNode.ID]
		var newFiber *Fiber

		if exists {
			delete(existingChildren, newNode.ID)

			newFiber = CreateWorkInProgress(matchedFiber, newNode.Props)

			if matchedFiber.Index != i {
				newFiber.Flags |= EffectMove
			}
			newFiber.Index = i

			if matchedFiber.HasPropsChanged() {
				newFiber.MarkUpdate()
			}
		} else {
			// create new fiber if new child doesn't exist in old fiber
			newFiber = &Fiber{
				Node:  newNode,
				Tag:   protocol.GetTagFromStr(newNode.Type),
				Key:   newNode.ID,
				Flags: EffectPlacement, // Đánh dấu là Mới
			}
		}

		newFiber.Parent = returnFiber
		newFiber.Sibling = nil

		if i == 0 {
			resultingFirstChild = newFiber
		} else {
			prevSibling.Sibling = newFiber
		}
		prevSibling = newFiber
	}

	for _, childToDelete := range existingChildren {
		childToDelete.MarkDeletion()
		returnFiber.Deletions = append(returnFiber.Deletions, childToDelete)
	}

	returnFiber.Child = resultingFirstChild
}

func (r *Reconciler) completeWork(fiber *Fiber) {
	fiber.BubbleFlags()
}

func (r *Reconciler) CommitRoot() {
	finishedWork := r.WipRoot

	if finishedWork == nil {
		return
	}

	if finishedWork.SubtreeFlags != EffectNone || finishedWork.Flags != EffectNone {
		r.commitWork(finishedWork)
	}

	r.CurrentRoot = finishedWork
	r.WipRoot = nil

	// --- GIAI ĐOẠN 3: POST-COMMIT (Render) ---
	// Bây giờ dữ liệu đã sạch sẽ, Scene Graph đã update, R-Tree đã update.
	// Ta báo cho Engine biết để sinh OpCode vẽ ra màn hình.

	// Lưu ý: Việc Render ra OpCode buffer thường được gọi bên ngoài (từ Engine)
	// sau khi WorkLoop trả về 'false' (đã xong).
	// Nhưng bạn cũng có thể bắn event từ đây nếu muốn.
}

func (r *Reconciler) commitWork(fiber *Fiber) {
	if fiber == nil {
		return
	}

	if len(fiber.Deletions) > 0 {
		for _, childToDelete := range fiber.Deletions {
			r.commitDeletion(childToDelete)
		}

		fiber.Deletions = nil
	}

	if fiber.Flags&EffectPlacement != 0 {
		// TODO: Complete this
	}
	if fiber.Flags&EffectUpdate != 0 {
		// TODO: Complete this
	}

	if fiber.SubtreeFlags != EffectNone {
		r.commitWork(fiber.Child)
		r.commitWork(fiber.Sibling)
	}
}

func (r *Reconciler) commitDeletion(fiberToDelete *Fiber) {
	if fiberToDelete == nil {
		return
	}

	if fiberToDelete.Node != nil {
		eng := r.engine.(*engine.Engine)
		eng.Spatial.Delete(fiberToDelete.Node.LastWorldMBR, fiberToDelete.Node.ID)

		if fiberToDelete.Node.Type == "INPUT" {
			eng.RemoveDOMOverlay(fiberToDelete.Node.ID)
		}
	}

	child := fiberToDelete.Child
	for child != nil {
		r.commitDeletion(child)
		child = child.Sibling
	}

	fiberToDelete.Alternate = nil
	fiberToDelete.Node = nil
}
