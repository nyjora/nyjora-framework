#!/usr/bin/env python
#-*- coding: utf-8 -*-

import sys, os
from optparse import OptionParser
import xml.etree.ElementTree as ET
from sets import Set
import jinja2

FALSE_VALUE = ["false", "False", "0"]
MAX_METHOD_TYPE = 8912

class ParseError(RuntimeError):
	pass

class Interface(object):
	def __init__(self):
		self.name = ""
		self.methods = []
		self.method_types = Set()
		self.method_names = Set()
		self.go_package = ""
		self.go_proto_package = ""
		self.go_private_imports = []
	
	def __str__(self):
		buf = "interface: name=%s\n" % self.name
		for m in self.methods:
			buf += str(m)
		return buf

class Method(object):
	def __init__(self):
		self.type = -1
		self.is_request = False
		self.name = ""
		self.arg = ""
		self.res = ""
		self.reliable = True
		self.async = False
		self.timeout = 0

	def __str__(self):
		return "\t%s: name=%s, type=%d, arg=%s, reliable=%s, async=%s\n" % \
				("request" if self.is_request else "message", self.name, self.type, \
				 self.arg, self.reliable, self.async)

class InterfaceParser(object):
	def __init__(self, verbose):
		self._tree = None
		self._interfaces = []
		self._verbose = verbose
		self._go_package = ""
		self._go_proto_package = ""
		self._go_imports = []
	
	def getResult(self):
		return self._interfaces
	
	def printResult(self):
		for interface in self._interfaces:
			print interface

	def parse(self, filename):
		tree = ET.parse(filename)
		self._go_package = tree.getroot().attrib["go.package"]
		self._go_proto_package = tree.getroot().attrib["go.proto_package"]
		for i in tree.getroot().iterfind("go.import"):
			self._go_imports.append(i.attrib["pkg"])

		for intNode in tree.getroot().iterfind("interface"):
			self._interfaces.append(self.parseInterface(intNode))

		if self._verbose:
			self.printResult();

	def parseInterface(self, intNode):
		interface = Interface()
		self._cur_interface = interface
		interface.name = intNode.attrib["name"]
		interface.go_package = self._go_package
		interface.go_imports = self._go_imports
		interface.go_proto_package = self._go_proto_package
		for i in intNode.iterfind("go.import"):
			interface.go_private_imports.append(i.attrib["pkg"])
		
		for metNode in intNode.iter():
			# parse methods
			method = self.parseMethod(metNode)
			if method is not None:
				interface.methods.append(method)

		return interface
	
	def parseMethod(self, metNode):
		if metNode.tag not in ["message", "request"]:
			return None

		method = Method()
		method.is_request = (metNode.tag == "request")
		try:
			method.name = metNode.attrib["name"]
			method.arg = metNode.attrib["arg"]
			method.type = int(metNode.attrib["type"])
	
			if method.type in self._cur_interface.method_types:
				raise ParseError("duplicated method type %d" % method.type)
			self._cur_interface.method_types.add(method.type)
	
			if method.name in self._cur_interface.method_names:
				raise ParseError("duplicated method name %s" % method.name)
			self._cur_interface.method_names.add(method.name)
	
			if method.type <= 0 or method.type > MAX_METHOD_TYPE:
				raise ParseError("method type must be between 1 and %d: " % MAX_METHOD_TYPE)
	
			if method.is_request:
				method.res = metNode.attrib["res"]
	
				try:
					method.timeout = int(metNode.attrib["timeout"])
				except ValueError:
					raise ParseError("method %s timeout attrbute must be interger" % method.name)
		
				if method.timeout <= 0:
					raise ParseError("method %s timeout must be larger than 0: " % method.name)
	
		except KeyError as err:
			raise ParseError("%s nodes must have a '%s' attribute" % (metNode.tag, err.args[0]))
		
		if "reliable" in metNode.attrib:
			method.reliable = (metNode.attrib["reliable"] not in FALSE_VALUE)

		if "async" in metNode.attrib:
			method.async = (metNode.attrib["async"] not in FALSE_VALUE)

		return method

class InterfaceCompilerGo:
	def __init__(self):
		self._env = jinja2.Environment(
			loader = jinja2.FileSystemLoader(os.path.dirname(__file__) + "/template"),
			keep_trailing_newline = True, # 文件结尾有换行符
			lstrip_blocks = True, # 去掉{% %}前面的空格和tab
			trim_blocks = True) # 去掉{% %}行符
		self._env.filters.update({
			"capfirst": self._capfirst
		})
		self._template_map = {
			"interface.go": "interface.go.jinja2",
		}

	def compile(self, interfaces, directory):
		for interface in interfaces:
			self.compileTemplateFor(directory, interface)
	
	def compileTemplateFor(self, directory, interface):
		for (source, templatename) in self._template_map.iteritems():
			templ = self._env.get_template(templatename)
			content = templ.render(interface.__dict__)
			f = open(directory + os.sep + interface.name.lower() + "_" + source, "w+")
			f.write(content.encode("UTF-8"))
	
	@staticmethod
	def _capfirst(inputstr):
		if len(inputstr) == 0:
			return ""

		return inputstr[0].upper() + inputstr[1:]

def main():
	optparser = OptionParser()
	optparser.add_option("--go_out", dest="go_out",
			help="directory for Go files", metavar="DIR")
	optparser.add_option("-v", "--verbose", dest="verbose",
			help="print a lot", action="store_true", default=False)

	optparser.usage = "%prog <path-to-interface.xml> [options]"

	(options, args) = optparser.parse_args()

	parser = InterfaceParser(options.verbose)
	for f in args:
		parser.parse(f)
	
	if options.go_out:
		compiler = InterfaceCompilerGo()
		compiler.compile(parser.getResult(), options.go_out)

if __name__ == "__main__":
	main()
