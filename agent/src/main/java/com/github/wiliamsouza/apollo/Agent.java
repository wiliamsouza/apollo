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
		String config = null;
		CommandLine cmd;

		Option configFile  = OptionBuilder.withArgName("file")
						  .hasArg()
						  .withDescription("Apollo agent configuration file.")
						  .create("config");
		Options options = new Options();
		options.addOption(configFile);
		options.addOption("h", "help", false, "Print this message.");

		CommandLineParser parser = new BasicParser();
		try {
			cmd = parser.parse(options, args);
		 	config = cmd.getOptionValue("config");
			if (cmd.hasOption("help")) {
				HelpFormatter formatter = new HelpFormatter();
				formatter.printHelp("apollo", options);
			}
		}
		catch (ParseException e){
			System.err.println("Option error: " + e.getMessage());
		}

		if (config == null) {
			config = "/etc/apollo/agent.conf";
		}

		System.out.println(config);
		System.out.println("Apollo agent. \n");
	}
}
