var prefix = "~uname_test"

const Discord = require('discord.js');
const client = new Discord.Client();
const settings = require('./settings.json');

client.on('ready',() => {
	console.log('Bot started');
});

client.on('message', message => {
	if (message.author === client.user) return;
	if (message.content.startsWith(prefix)) {
		console.log('Got a message for me: '+message.content);
		var cmnd_array = message.content.split(" ");
		message.guild.channels.forEach(function(channel){
			if ( channel.name === cmnd_array[1] && channel.type === "voice" ){
				var yes = [];
                                var no = [];
				var no2 = [];

				channel.members.forEach(function(member){
					name = "";
					found = false;
					if (member.nickname == null) {
						name = member.user.username;
					}else{
						name = member.nickname;
					}
					cmnd_array.slice(2).forEach(function(name_g){
						if (name_g === name){
							found = true;
						}
					})
					if(found){
						yes.push(name);
					}else{
						no.push(name);
					}
				})
				cmnd_array.slice(2).forEach(function(name_g){
					found = false;
					yes.forEach(function(name_on){
						if (name_g === name_on){
                                                        found = true;
                                                }
					})
					if(found){
						
					}else{
						no2.push(name_g);
					}
				})
				message.channel.sendMessage('Playes found    : '+yes);
				message.channel.sendMessage('Unknown Discord nicknames: '+no);
				message.channel.sendMessage('Unknown Ingame players: '+no2);
			}
		})
	}
});

client.login(settings.token);
