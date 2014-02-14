package com.github.wiliamsouza.apollo;

import java.io.IOException;
import java.net.URI;

import org.apache.commons.cli.Options;
import org.apache.commons.cli.Option;
import org.apache.commons.cli.OptionBuilder;
import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.CommandLineParser;
import org.apache.commons.cli.BasicParser;
import org.apache.commons.cli.HelpFormatter;
import org.apache.commons.cli.ParseException;

import javax.websocket.ContainerProvider;
import javax.websocket.DeploymentException;
import javax.websocket.Session;
import javax.websocket.WebSocketContainer;

public class Agent {

    public static void main(String[] args) {

        String config = "/etc/apollo/agent.conf";
        CommandLine cmd;
        Session session = null;

        Option configFile = OptionBuilder.withArgName("file")
                          .hasArg()
                          .withDescription("Apollo agent configuration file.")
                          .create("config");
        Options options = new Options();
        options.addOption(configFile);
        options.addOption("h", "help", false, "Print this message.");

        CommandLineParser parser = new BasicParser();
        try {
            cmd = parser.parse(options, args);
            if (cmd.hasOption("help")) {
                HelpFormatter formatter = new HelpFormatter();
                formatter.printHelp("apollo", options);
            }
            String conf = cmd.getOptionValue("config");
            if (conf != null) {
                config = conf;
            }
        }
        catch (ParseException e){
            System.err.println("Option error: " + e.getMessage());
        }

        // TODO: Add the following as options to etc/apollo/agent.conf file
        String APIKey = "";
        String ADBPath = "/usr/bin/adb";
        String serverURI = "ws://localhost:8000/ws/agent/";
        URI uri = URI.create(serverURI + APIKey);

        WebSocketContainer container = ContainerProvider.getWebSocketContainer();
        try {
            session = container.connectToServer(WebSocketEndpoint.class, uri);
        } catch (DeploymentException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        }

        DeviceMonitor monitor = new DeviceMonitor(session, ADBPath);
        monitor.start();
        //monitor.finish();

        System.out.println(config);
        System.out.println("Apollo agent. \n");

        while (true) {
        }
    }
}
