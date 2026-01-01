#include "kinect_bridge.h"

#include <stdlib.h>
#include <math.h>
#include <time.h>

// Mock implementation

void init_kinect() {
    srand(time(NULL));
}

int get_skeleton(KSkeleton* skeleton) {
    // Generate dummy data
    skeleton->tracked = 1;

    // Simulate some movement
    float t = (float)clock() / CLOCKS_PER_SEC;

    for (int i = 0; i < 25; i++) {
        skeleton->joints[i].x = 0;
        skeleton->joints[i].y = 0;
        skeleton->joints[i].z = 2.0; // 2 meters away
        skeleton->joints[i].state = HAND_STATE_NOT_TRACKED;
    }

    // Set specific joints used in Main.java
    // SpineBase, SpineMid, HandRight, HandLeft, SpineShoulder, ShoulderRight

    skeleton->joints[JOINT_SPINE_BASE].x = 0;
    skeleton->joints[JOINT_SPINE_BASE].y = 0;
    skeleton->joints[JOINT_SPINE_BASE].z = 2.0;

    skeleton->joints[JOINT_SPINE_MID].x = 0;
    skeleton->joints[JOINT_SPINE_MID].y = 0.3;
    skeleton->joints[JOINT_SPINE_MID].z = 2.0;

    skeleton->joints[JOINT_SPINE_SHOULDER].x = 0;
    skeleton->joints[JOINT_SPINE_SHOULDER].y = 0.5;
    skeleton->joints[JOINT_SPINE_SHOULDER].z = 2.0;

    skeleton->joints[JOINT_SHOULDER_RIGHT].x = 0.2;
    skeleton->joints[JOINT_SHOULDER_RIGHT].y = 0.5;
    skeleton->joints[JOINT_SHOULDER_RIGHT].z = 2.0;

    // Move right hand in a circle
    skeleton->joints[JOINT_HAND_RIGHT].x = 0.2 + 0.3 * cos(t);
    skeleton->joints[JOINT_HAND_RIGHT].y = 0.5 + 0.3 * sin(t);
    skeleton->joints[JOINT_HAND_RIGHT].z = 1.5;
    skeleton->joints[JOINT_HAND_RIGHT].state = HAND_STATE_OPEN; // Force creation condition

    skeleton->joints[JOINT_HAND_LEFT].x = -0.2;
    skeleton->joints[JOINT_HAND_LEFT].y = 0.5;
    skeleton->joints[JOINT_HAND_LEFT].z = 2.0;

    return 1;
}
