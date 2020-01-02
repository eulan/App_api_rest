var app = new Vue({
  el: '#app',
  data: {
    domain: '',
    message: '',
    endpoints: [],
    servers_changed: '',
    ssl_grade: '',
    previous_ssl_grade: '',
    logo: '',
    title: '',
    is_down: '',
    see: false
  },
  methods:{
  	traer: function(){
      var url = new URL('http://localhost:8000/='+this.domain)
      fetch(url)
        .then(data => data.json())
        .then(data => {
          console.log(data)
          this.endpoints = data.endpoints
          this.servers_changed = data.servers_changed
          this.ssl_grade = data.ssl_grade
          this.previous_ssl_grade = data.previous_ssl_grade
          this.logo = data.logo
          this.title = data.title
          this.is_down = data.is_down
          this.see = true
        })
  	}
  }
})