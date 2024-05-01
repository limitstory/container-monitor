package main

import (
	"context"
	"fmt"
	mod "monitor/modules"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"os"
)

func IsSucceed(podsItems []v1.Pod) bool {
	for _, pod := range podsItems {
		if pod.Status.Phase != "Succeeded" {
			return false
		}
	}
	return true
}

func main() {

	// kubernetes api 클라이언트 생성하는 모듈
	clientset := mod.InitClient()
	if clientset == nil {
		fmt.Println("Could not create client!")
		os.Exit(-1)
	}

	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	if IsSucceed(pods.Items) == false {
		time.Sleep(time.Second)
	} else {
		var startedTestTime int64 = 9999999999999
		var finishedTestTime int64 = 0
		var minContainerRunningTime int64 = 9999999999999
		var maxContainerRunningTime int64 = 0
		var totalContinerRunningTime int64 = 0
		var totalContainerRestart int64 = 0
		var ContainerRestartArray [20]int64

		for _, pod := range pods.Items {
			fmt.Println()
			fmt.Println(pod.Name)
			fmt.Println("PodContitions")
			for i := 0; i < len(pod.Status.Conditions); i++ {
				fmt.Println(pod.Status.Conditions[i].Type, ":", pod.Status.Conditions[i].LastTransitionTime.Unix())
			}
			startTime := pod.Status.StartTime.Unix()
			startedAt := pod.Status.ContainerStatuses[0].State.Terminated.StartedAt.Unix()
			finishedAt := pod.Status.ContainerStatuses[0].State.Terminated.FinishedAt.Unix()
			runningTime := finishedAt - startedAt

			fmt.Println("Time Info")
			fmt.Println("StartTime:", startTime)
			fmt.Println("StartedAt:", startedAt)
			fmt.Println("FinishedAt:", finishedAt)

			fmt.Println("running Time:", runningTime)
			fmt.Println("RestartCount:", pod.Status.ContainerStatuses[0].RestartCount)

			if startedTestTime > startTime {
				startedTestTime = startTime
			}
			if finishedTestTime < finishedAt {
				finishedTestTime = finishedAt
			}

			if minContainerRunningTime > runningTime {
				minContainerRunningTime = runningTime
			}
			if maxContainerRunningTime < runningTime {
				maxContainerRunningTime = runningTime
			}

			totalContinerRunningTime += runningTime
			totalContainerRestart += int64(pod.Status.ContainerStatuses[0].RestartCount)
			ContainerRestartArray[pod.Status.ContainerStatuses[0].RestartCount]++
			//a := pod.Status.
		}
		fmt.Println()
		fmt.Println("Total Info")
		fmt.Println("StartedTestTime:", startedTestTime)
		fmt.Println("FinishedTestTime", finishedTestTime)
		fmt.Println("TestingTime:", finishedTestTime-startedTestTime)
		fmt.Println("minContainerRunningTime:", minContainerRunningTime)
		fmt.Println("averageContainerRunningTime:", float64(totalContinerRunningTime)/float64(len(pods.Items)))
		fmt.Println("maxContainerRunningTime:", maxContainerRunningTime)
		fmt.Println("TotalContainerRestart", totalContainerRestart)
		fmt.Println("ContainerRestartArray", ContainerRestartArray)
	}
}
