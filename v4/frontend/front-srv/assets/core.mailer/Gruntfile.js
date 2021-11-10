module.exports = function(grunt) {

    const {initConfig, loadNpmTasks, registerTasks} = require('../gruntConfigCommon.js')
    grunt.initConfig(initConfig('PydioMailer'));
    loadNpmTasks(grunt);
    registerTasks(grunt);

};