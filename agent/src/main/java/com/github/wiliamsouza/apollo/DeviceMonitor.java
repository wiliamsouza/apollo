package com.github.wiliamsouza.apollo;

import com.android.ddmlib.AndroidDebugBridge;

public class DeviceMonitor {

    public void start() {
        AndroidDebugBridge.init(false);
        AndroidDebugBridge.addDeviceChangeListener(new DeviceListener());
        AndroidDebugBridge adb = AndroidDebugBridge.createBridge("adb", true);
        if (!adb.isConnected()) {
            System.out.println("Couldn't connect to ADB server");
        }
    }

    public void finish() {
        AndroidDebugBridge.disconnectBridge();
        AndroidDebugBridge.terminate();
    }
}
