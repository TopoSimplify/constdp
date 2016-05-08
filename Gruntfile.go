package constdp
module.exports = func (grunt) {
  grunt.initConfig({
    tape : {
      options: {
        pretty: true,
        output: 'console'
      },
      files  : ['test/*.js']
    },
    watch: {
      scripts: {
        files  : ['test/*/*.js','lib/*/*.js'],
        tasks  : 'tape',
        options: {
          debounceDelay: 200
        }
      }
    }
  })

  grunt.loadNpmTasks('grunt-contrib-watch')
  grunt.loadNpmTasks('grunt-tape')
  grunt.registerTask('test', ['tape'])
  grunt.registerTask('default', ['test'])

  grunt.task.run(['test','watch'])

}
