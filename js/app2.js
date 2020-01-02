var app2 = new Vue({
  el: '#app2',
  data: {
    items: [],
    watch: false
  },
  methods:{
    watching: function () {
      this.watch = true
    }
  },
  created(){
    var url = new URL('http://localhost:8000/consults')
      fetch(url)
        .then(data => data.json())
        .then(data => {
          console.log(data)
          this.items = data.items
        })
  }
})