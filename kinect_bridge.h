#ifndef KINECT_BRIDGE_H
#define KINECT_BRIDGE_H

typedef struct {
    float x;
    float y;
    float z;
    int state;
} KJoint;

typedef struct {
    KJoint joints[25];
    int tracked;
} KSkeleton;

// Joint Types (Standard Kinect v2)
#define JOINT_SPINE_BASE 0
#define JOINT_SPINE_MID 1
#define JOINT_NECK 2
#define JOINT_HEAD 3
#define JOINT_SHOULDER_LEFT 4
#define JOINT_ELBOW_LEFT 5
#define JOINT_WRIST_LEFT 6
#define JOINT_HAND_LEFT 7
#define JOINT_SHOULDER_RIGHT 8
#define JOINT_ELBOW_RIGHT 9
#define JOINT_WRIST_RIGHT 10
#define JOINT_HAND_RIGHT 11
#define JOINT_HIP_LEFT 12
#define JOINT_KNEE_LEFT 13
#define JOINT_ANKLE_LEFT 14
#define JOINT_FOOT_LEFT 15
#define JOINT_HIP_RIGHT 16
#define JOINT_KNEE_RIGHT 17
#define JOINT_ANKLE_RIGHT 18
#define JOINT_FOOT_RIGHT 19
#define JOINT_SPINE_SHOULDER 20
#define JOINT_HAND_TIP_LEFT 21
#define JOINT_THUMB_LEFT 22
#define JOINT_HAND_TIP_RIGHT 23
#define JOINT_THUMB_RIGHT 24

// Hand States
#define HAND_STATE_UNKNOWN 0
#define HAND_STATE_NOT_TRACKED 1
#define HAND_STATE_OPEN 2
#define HAND_STATE_CLOSED 3
#define HAND_STATE_LASSO 4

void init_kinect();
int get_skeleton(KSkeleton* skeleton);

#endif
