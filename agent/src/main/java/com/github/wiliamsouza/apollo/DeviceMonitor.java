package com.github.wiliamsouza.apollo;

import com.android.ddmlib.AndroidDebugBridge;
import com.android.ddmlib.IDevice;


import com.github.wiliamsouza.apollo.DeviceListener;

public class DeviceMonitor {

    public void start() {
        AndroidDebugBridge.init(false);
    }

    public void finish() {
        AndroidDebugBridge.terminate();
    }

    public void run() throws Exception {
        AndroidDebugBridge.addDeviceChangeListener(new DeviceListener());

        AndroidDebugBridge adb = AndroidDebugBridge.createBridge("adb", true);

        Thread.sleep(1000);
        if (!adb.isConnected()) {
            System.out.println("Couldn't connect to ADB server");
        }

        //AndroidDebugBridge.disconnectBridge();
    }
}
