// Copyright © 2019. All rights reserved.
// Author: Ilya Yuryevich.
// Contacts: qioalice@gmail.com, https://github.com/qioalice
// License: https://opensource.org/licenses/MIT

package log

// Logger is used to generate and write log messages.
//
// You can instantiate as many your own loggers with different behaviour,
// different contexts, as you want. But also you can just use package level logger,
// modernize and configure it the same way as any instantiated Logger object.
//
// Inheritance.
//
// Remember!
// No one func or method does not change the current object.
// They always creates and returns a copy of the current object with applied
// your changes except in cases in which the opposite is not explicitly indicated.
//
// Of course you can chain all setters and then call log message generator, like this:
//
// 		log := Package("main")
// 		log.With("key", "value").Warn("It's dangerous!")
//
// But in that case, after finishing execution of second line,
// 'log' variable won't contain add field "key" and group "group".
// But package name "main" will contain.
// Want a different behaviour and want to have Logger with these fields?
// No problem, save generated Logger object:
//
// 		log := Package("main")
// 		log = log.Group("group").With("key", "value")
// 		log.Warn("It's dangerous!")
//
// Because of all finishers (methods that actually writes a log message, e.g:
// Debug, Debugf, Debugw, Warn, Warnf, Warnw, etc...) also returns a Logger
// object they uses to generate log Entry, you can save it too, and finally
// it's the same as in the example above:
//
// 		log := Package("main")
// 		log = log.Group("group").With("key", "value").Warn("It's dangerous!")
//
// but it's strongly not recommended to do so, because it made code less clear.
//
// There are 5 Logger constructors:
//
// 		Package(packageName, options...)
// 		Func(funcName, options...)
// 		Class(className, options...)
// 		Method(className, methodName, options...)
// 		New(options...)
//
// You can instantiate Logger object using any of constructor listed above.
// First four are used to create Logger object that binds to some Golang entity,
// and their output will contain field with 'sys.func' key and your passed value.
//
// You can combine them then to explicitly create an exactly Logger you want:
// E.g.: Method(className, methodName) == Class(className).Method(methodName).
// In that case 'sys.func' will have this value: 'className.methodName'.
//
// And the fifth creates a common regular Logger object, but it will contain
// 'sys.func' field also. Because there is auto-generation reflect information
// by default (based by stacktrace).
// You can disable this behaviour applying 'Options.EnableCallerInfo(false)'
// option to Logger's constructor or using 'Apply' method.
type Logger struct {

	// Core copies much less times than Entry.
	// Core contains log messages destinations, auxiliary functions, log format,
	// and something like this.
	core *core

	// Entry is what log message is.
	// Entry is it's stacktrace, caller info, timestamp, level, message, group,
	// flags, etc.
	entry *Entry
}

// Apply overwrites current Logger's copy behaviour by provided reasons.
// Read more: Options.
func (l *Logger) Apply(options ...interface{}) (copy *Logger) {

	if len(options) == 0 || !l.canContinue() {
		return l
	}
	return l.apply(options)
}

// Package sets current Logger's copy package name.
// Zeroes Logger's function name or Logger's class and method names.
// Disables caller auto-generation.
func (l *Logger) Package(packageName string) (copy *Logger) {

	if packageName == "" || !l.canContinue() {
		return l
	}
	return l.derive(false).entry.setPackageName(packageName).l
}

// Func sets current Logger's copy function name.
// Zeroes Logger's class and method names. Disables auto generating caller.
func (l *Logger) Func(funcName string) (copy *Logger) {

	if funcName == "" || !l.canContinue() {
		return l
	}
	return l.derive(false).entry.setFuncName(funcName).l
}

// Class sets current Logger's copy class name.
// Zeroes logger's function and method names. Disables auto generating caller.
func (l *Logger) Class(className string) (copy *Logger) {

	if className == "" || !l.canContinue() {
		return l
	}
	return l.derive(false).entry.setClassName(className).l
}

// Method sets current Logger's copy method name.
// Zeroes logger's function name. Disables auto generating caller.
//
// If class name isn't set, there is the same behaviour as Func.
func (l *Logger) Method(methodName string) (copy *Logger) {

	if methodName == "" || !l.canContinue() {
		return l
	}
	return l.derive(false).entry.setMethodName(methodName).l
}

// With adds fields to the current Logger's copy.
//
// You can pass both of explicit or implicit fields. Even both of named/unnamed
// implicit fields, but names (keys) should be only string.
// Neither string-like (fmt.Stringer) nor string-cast ([]byte). Only strings.
func (l *Logger) With(fields ...interface{}) (copy *Logger) {

	if len(fields) == 0 || !l.canContinue() {
		return l
	}
	return l.derive(false).entry.with(fields, nil).l
}

// WithStrict adds an explicit fields to the current Logger's copy.
func (l *Logger) WithStrict(fields ...Field) (copy *Logger) {

	if len(fields) == 0 || !l.canContinue() {
		return l
	}
	return l.derive(false).entry.with(nil, fields).l
}