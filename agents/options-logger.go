package agents

func WithLogger(logger *Logger) AgentOption {
	return func(agent *Agent) {
		agent.logger = logger
	}
}

func WithLogging(level LogLevel, enabled bool) AgentOption {
	return func(agent *Agent) {
		agent.logger = NewLogger(level, enabled)
	}
}

func WithLoggingEnabled() AgentOption {
	return func(agent *Agent) {
		if agent.logger == nil {
			agent.logger = NewLogger(LogLevelInfo, true)
		} else {
			agent.logger.SetEnabled(true)
		}
	}
}

func WithLoggingDisabled() AgentOption {
	return func(agent *Agent) {
		if agent.logger == nil {
			agent.logger = NewLogger(LogLevelOff, false)
		} else {
			agent.logger.SetEnabled(false)
		}
	}
}

func WithLogLevel(level LogLevel) AgentOption {
	return func(agent *Agent) {
		if agent.logger == nil {
			agent.logger = NewLogger(level, true)
		} else {
			agent.logger.SetLevel(level)
		}
	}
}