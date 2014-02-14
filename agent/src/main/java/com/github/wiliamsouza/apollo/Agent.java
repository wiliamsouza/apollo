package com.github.wiliamsouza.apollo;

import org.apache.commons.cli.Options;
import org.apache.commons.cli.Option;
import org.apache.commons.cli.OptionBuilder;
import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.CommandLineParser;
import org.apache.commons.cli.BasicParser;
import org.apache.commons.cli.HelpFormatter;
import org.apache.commons.cli.ParseException;

public class Agent {

    public static void main(String[] args) {
        String config = "/etc/apollo/agent.conf";
        CommandLine cmd;

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

        DeviceMonitor monitor = new DeviceMonitor();
        monitor.start();
        //monitor.finish();

        System.out.println(config);
        System.out.println("Apollo agent. \n");

        while (true) {
        }
    }
}
