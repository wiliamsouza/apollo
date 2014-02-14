package com.github.wiliamsouza.apollo;

import com.android.ddmlib.AdbCommandRejectedException;
import com.android.ddmlib.IDevice;
import com.android.ddmlib.AndroidDebugBridge.IDeviceChangeListener;
import com.android.ddmlib.TimeoutException;

import javax.websocket.Session;
import java.io.IOException;

public class DeviceListener implements IDeviceChangeListener {

    private Session session;

    public DeviceListener(Session session) {
        this.session = session;
    }

    public void deviceConnected(IDevice device) {
        String msg = "Device connected " + device.getSerialNumber();
        System.out.println(msg);
        System.out.println(device.getState());
        try {
            this.session.getBasicRemote().sendText(msg);
        } catch (IOException e) {
            e.printStackTrace();
        }
        System.out.println("**********************************************************************");
    }

    public void deviceDisconnected(IDevice device) {
        String msg = "Device disconnected " + device.getSerialNumber();
        System.out.println(msg);
        System.out.println(device.getState());
        try {
            this.session.getBasicRemote().sendText(msg);
        } catch (IOException e) {
            e.printStackTrace();
        }
        System.out.println("**********************************************************************");
    }

    public void deviceChanged(IDevice device, int changeMask) {
        String msg = "Device changed " + device.getSerialNumber();
        System.out.println(msg);
        System.out.println(device.getState());
        try {
            this.session.getBasicRemote().sendText(msg);
        } catch (IOException e) {
            e.printStackTrace();
        }
        System.out.println("**********************************************************************");
    }
}