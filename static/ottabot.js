var userName = null;
const testLogin = Vue.extend({
	template:"#logoutButton",
	data:function(){
		return {userName:userName};
	},
	computed:{
		userName:function(){
			if (this.userName == null){
				return "Login";
			}else{
				return this.userName;
			}
		}
	}
});
var paths = [{path:"/",component:testLogin}];
var Router = new VueRouter({routes:paths});
var App = new Vue({
	router:Router,
	el:"#app"
});
//Router.start(App,"#app");