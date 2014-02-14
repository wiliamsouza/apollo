package com.github.wiliamsouza.apollo;

import com.android.ddmlib.AndroidDebugBridge;

import javax.websocket.Session;

public class DeviceMonitor {

    private Session session;
    private String ADBPath;

    public DeviceMonitor(Session session, String ADBPath) {
        this.session = session;
        this.ADBPath = ADBPath;
    }
    public void start() {
        AndroidDebugBridge.init(false);
        AndroidDebugBridge.addDeviceChangeListener(new DeviceListener(session));
        AndroidDebugBridge adb = AndroidDebugBridge.createBridge(ADBPath, true);

        try {
            Thread.sleep(1000);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        if (!adb.isConnected()) {
            System.out.println("Couldn't connect to ADB server");
        }
    }

    public void finish() {
        AndroidDebugBridge.disconnectBridge();
        AndroidDebugBridge.terminate();
    }
}
