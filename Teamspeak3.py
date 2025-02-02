#!/usr/bin/env python
# -*- encoding: utf-8; py-indent-offset: 4 -*-

# This is free software;  you can redistribute it and/or modify it
# under the  terms of the  GNU General Public License  as published by
# the Free Software Foundation in version 2.  This file is distributed
# in the hope that it will be useful, but WITHOUT ANY WARRANTY;  with-
# out even the implied warranty of  MERCHANTABILITY  or  FITNESS FOR A
# PARTICULAR PURPOSE. See the  GNU General Public License for more de-
# ails.  You should have  received  a copy of the  GNU  General Public
# License along with GNU Make; see the file  COPYING.  If  not,  write
# to the Free Software Foundation, Inc., 51 Franklin St,  Fifth Floor,
# Boston, MA 02110-1301 USA.

# <<<Teamspeak3>>>
# ConfigError: No
# QueryPortReachable: Yes
# AuthSuccess: Yes
# Version: 3.12.1
# Platform: Linux
# Build: 1585305527
# VirtualServer: (9987 online 0 32 0 yes 0 0)

def parse_teamspeak3(info):
    parsed = {u'VirtualServer': {}}
    for line in info:
        if line[0][:-1] == "VirtualServer":
            data = {u'port': line[1][1:],
                    u'status': line[2],
                    u'clientsonline': int(line[3]),
                    u'clientsmax': int(line[4]),
                    u'channels': int(line[5]),
                    u'autostart': line[6],
                    u'ingress': int(line[7]),
                    u'egress': int(line[8][:-1])}
            parsed[u'VirtualServer'][data['port']] = data
        else:
            parsed[line[0][:-1]] = " ".join(line[1:])
    return parsed


def inventory_teamspeak3(parsed):
    if parsed.get(u'Version'):
        yield u'Global', {}
    for port, vs in parsed['VirtualServer'].items():
        yield port, {u'status': vs[u'status'], u'autostart': vs[u'autostart']}


def check_teamspeak3(item, params, parsed):
    if 'ConfigError' in parsed:
        if parsed['ConfigError'] == 'Yes, 1':
            return 1, "Unable to read the agent plugin configuration from teamspeak3.cfg"
        elif parsed['ConfigError'] == 'Yes, 2':
            return 1, "No section serverquery in teamspeak3.cfg"
        elif parsed['ConfigError'] == 'Yes, 3':
            return 1, "No Teamspeak3 server address in teamspeak3.cfg"
        elif parsed['ConfigError'] == 'Yes, 4':
            return 1, "No Teamspeak3 server user in teamspeak3.cfg"
        elif parsed['ConfigError'] == 'Yes, 5':
            return 1, "No Teamspeak3 server password in teamspeak3.cfg"
    if parsed.get('QueryPortReachable', 'No') == 'No':
        return 2, "Server unreachable"
    if parsed.get('AuthSuccess', 'No') == 'No':
        return 2, "Unable to authenticate"
    if item in parsed[u'VirtualServer']:
        vs = parsed[u'VirtualServer'][item]
        state = 0
        text = vs[u'status']
        if vs[u'status'] != params[u'status']:
            text += '(!!)'
            state = max(state, 2)
        text += " (%d clients, %d channels), autostart %s" % (vs[u'clientsonline'], vs[u'channels'], vs[u'autostart'])
        if vs[u'autostart'] != params[u'autostart']:
            text += '(!)'
            state = max(state, 1)
        now = time.time()
        in_rate = get_rate('teamspeak3.%d.if_in_octets' % int(item), now, vs[u'ingress'])
        out_rate = get_rate('teamspeak3.%d.if_out_octets' % int(item), now, vs[u'egress'])
        perfdata = [(u'current_users', vs[u'clientsonline'], None, None, 0, vs[u'clientsmax']),
                    (u'channels', vs[u'channels']),
                    (u'if_in_octets', in_rate),
                    (u'if_out_octets', out_rate)]
        return state, text, perfdata
    if item == u'Global':
        return 0, "Platform: %s, Version: %s %s" % (parsed[u'Platform'], parsed[u'Version'], parsed[u'Build'])


check_info['Teamspeak3'] = {
    'parse_function': parse_teamspeak3,
    'inventory_function': inventory_teamspeak3,
    'check_function': check_teamspeak3,
    'service_description': "Teamspeak3 %s",
    'has_perfdata': True,
}
