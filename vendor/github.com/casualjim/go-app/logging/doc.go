/*Package logging provides a configuration model for logrus.

It introduces hierarchical configuration and hot-reloading of configuration for logrus through viper
This means that you can use YAML, TOML, HCL and JSON to configure your loggers, or use env vars.
It allows for storing configuration in etcd or consul.

You can configure multiple root level loggers which serve as default logger for that tree. In addition
to being the default logger, each child logger created will inherit the configuration from its parents.
So you only have to define overrides for things you want to customize.

    logging:
      root:
        level: debug
        hooks:
          - name: journald
        writer: stderr
        child1:
          level: info
          hooks:
            - name: file
              path: ./app.log
            - name: syslog
              network: udp
              host: localhost:514
              priority: info
      alerts:
        level: error
        writer: stderr

In this example config there are 2 root loggers defined: root and alerts.
Both root loggers use a stderr writer, other writers are stdout and discard.

The logger named root has also one child node that overrides the level property.
The hooks property is not an override but instead combines the defined hooks.
So the hooks work additive whereas everything else are overrides.

Custom writers:

You can register your own named writers with the RegisterWriter function.

Custom Formatters:

You can register your own formatters with the RegisterFormatter function.
*/
package logging
