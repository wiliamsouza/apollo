package com.github.wiliamsouza.apollo;

import com.android.ddmlib.IDevice;
import com.android.ddmlib.AndroidDebugBridge.IDeviceChangeListener;

public class DeviceListener implements IDeviceChangeListener {

    public void deviceConnected(IDevice device) {
        System.out.println("Device connected " + device.getSerialNumber());
    }

    public void deviceDisconnected(IDevice device) {
        System.out.println("device disconnected " + device.getSerialNumber());
    }

    public void deviceChanged(IDevice device, int changeMask) {
        System.out.println("Device changed " + device.getSerialNumber());
    }
}