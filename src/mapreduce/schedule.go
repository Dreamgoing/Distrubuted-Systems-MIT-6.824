package mapreduce

import (
	"fmt"
	"sync"
	"time"
)

//
// schedule() starts and waits for all tasks in the given phase (Map
// or Reduce). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//
	GolangChannelSolution(jobName, mapFiles, nReduce, phase, registerChan)
	fmt.Printf("Schedule: %v phase done\n", phase)
}

const (
	WorkerProcessing = iota
	WorkerIdle
)

type WorkerInfo struct {
	sync.Mutex

	address string
	state   int
}

type WorkerQueue struct {
	sync.RWMutex
	newCond *sync.Cond

	queue []*WorkerInfo
}

func newWorkerQueue() (wq *WorkerQueue) {
	wq = new(WorkerQueue)
	wq.newCond = sync.NewCond(wq)
	return
}

func (w *WorkerInfo) changeState(state int) {
	w.Lock()
	w.state = state
	w.Unlock()
}

func (w *WorkerInfo) isIdle() bool {
	return w.state == WorkerIdle
}

func (wq *WorkerQueue) ChooseIdleWorker() *WorkerInfo {
	//wq.Lock()
	//defer wq.Unlock()
	for _, it := range wq.queue {
		if it.isIdle() {
			return it
		}
	}
	return nil
}

func (wq *WorkerQueue) AddWorker(worker string) {
	wq.Lock()
	defer wq.Unlock()
	wq.queue = append(wq.queue, &WorkerInfo{address: worker, state: WorkerIdle})
	wq.newCond.Signal()
	fmt.Println("AddWorker waiting.")
	fmt.Println(len(wq.queue))
}

func CondMutexSolution(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		fmt.Println("Schedule map")
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		fmt.Println("Schedule reduce")
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)
	workerQueue := newWorkerQueue()

	go func() {
		for it := range registerChan {
			workerQueue.AddWorker(it)
		}
	}()

	/// 有点类似, 经典的生产者消费者模式.
	var wg sync.WaitGroup
	wg.Add(ntasks)
	for i := 0; i < ntasks; i++ {
		/// Call rpc on idle worker

		go func(taskNumber int) {
			rpcArgs := &DoTaskArgs{
				JobName:       jobName,
				File:          mapFiles[taskNumber],
				Phase:         phase, TaskNumber: taskNumber,
				NumOtherPhase: n_other}
			rpcRes := false
			worker := &WorkerInfo{}

			for {
				workerQueue.Lock()
				worker = workerQueue.ChooseIdleWorker()
				//fmt.Printf("Task %d: worker: %v\n", taskNumber, worker)
				if worker != nil {
					fmt.Println("address:", worker.address)
					worker.changeState(WorkerProcessing)
					workerQueue.Unlock()
					break
				} else {

					fmt.Printf("Task%d. Waiting\n", taskNumber)
					workerQueue.newCond.Wait()
				}
				workerQueue.Unlock()
			}

			rpcRes = call(worker.address, "Worker.DoTask", rpcArgs, nil)
			worker.changeState(WorkerIdle)
			workerQueue.newCond.Signal()

			wg.Done()
			fmt.Printf("Task %d/%d finish. rpcRes: %v", taskNumber, ntasks, rpcRes)
		}(i)

	}
	wg.Wait()

	fmt.Printf("Schedule: %v phase done\n", phase)
}

func CallWithRetry(workers chan string, rpcArgs *DoTaskArgs) {
	for {
		worker := <-workers
		rpcRes := call(worker, "Worker.DoTask", rpcArgs, nil)
		if rpcRes {
			workers <- worker
			break
		}

		// Delay 3s to put failed worker to worker queue.
		go func(failedWorker string) {
			time.Sleep(time.Second * 3)
			workers <- failedWorker
		}(worker)
	}
}

func GolangChannelSolution(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		fmt.Println("Schedule map")
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		fmt.Println("Schedule reduce")
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	workers := make(chan string, 3)

	go func() {
		for it := range registerChan {
			workers <- it
		}
	}()

	var wg sync.WaitGroup
	wg.Add(ntasks)
	for i := 0; i < ntasks; i++ {
		go func(taskNumber int) {
			rpcArgs := &DoTaskArgs{
				JobName:       jobName,
				File:          mapFiles[taskNumber],
				Phase:         phase, TaskNumber: taskNumber,
				NumOtherPhase: n_other}
			CallWithRetry(workers, rpcArgs)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
func showArgs(jobName string, mapFiles []string, nReduce int, phase jobPhase) {
	fmt.Println("[showArgs]")
	fmt.Println("jobName: ", jobName)
	fmt.Println("mapFiles: ")
	for _, it := range mapFiles {
		fmt.Print(it, " ")
	}
	fmt.Print("\n")
	fmt.Println("nReduce: ", nReduce)
	fmt.Println("phase: ", phase)
}
