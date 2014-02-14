package com.github.wiliamsouza.apollo;

import javax.websocket.ClientEndpoint;
import javax.websocket.OnMessage;
import javax.websocket.OnOpen;
import javax.websocket.OnClose;
import javax.websocket.OnError;
import javax.websocket.Session;
import javax.websocket.CloseReason;

@ClientEndpoint
public class WebSocketEndpoint {

    @OnOpen
    public void onOpen( final Session session ){
        System.out.println("Open session.");
    }

    @OnMessage
    public void onMessage(String message, Session session) {
        System.out.println("Message: " + message);
    }

    @OnClose
    public void onClose(Session session, CloseReason closeReason) {
        System.out.println("Close: " + closeReason.getReasonPhrase());
    }

    @OnError
    public void onError(Session session, Throwable thr) {
        System.out.println("Error: " + thr.getMessage());
    }
}