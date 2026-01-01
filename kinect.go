package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -lm
#include "kinect_bridge.h"
*/
import "C"

type KinectDevice struct {
	skeleton C.KSkeleton
}

func NewKinectDevice() *KinectDevice {
	C.init_kinect()
	return &KinectDevice{}
}

func (k *KinectDevice) Update() bool {
	res := C.get_skeleton(&k.skeleton)
	return res != 0
}

func (k *KinectDevice) GetJoint(jointType int) (float64, float64, float64, int) {
	j := k.skeleton.joints[jointType]
	return float64(j.x), float64(j.y), float64(j.z), int(j.state)
}

const (
	JointType_SpineBase     = C.JOINT_SPINE_BASE
	JointType_SpineMid      = C.JOINT_SPINE_MID
	JointType_HandRight     = C.JOINT_HAND_RIGHT
	JointType_HandLeft      = C.JOINT_HAND_LEFT
	JointType_SpineShoulder = C.JOINT_SPINE_SHOULDER
	JointType_ShoulderRight = C.JOINT_SHOULDER_RIGHT
	HandState_Open          = C.HAND_STATE_OPEN
)
