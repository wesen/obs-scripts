obs = obslua
source_def = {}
source_def.id = "button_example_script"

-- Function called when the button is clicked
function button_clicked(props, prop)
    obs.script_log(obs.LOG_INFO, "Button clicked!")
    return false
end

-- Define the script properties
function script_properties()
    local props = obs.obs_properties_create()
    local button = obs.obs_properties_add_button(props, "button", "Click Me!", button_clicked)
    return props
end

-- Register the script as a source
function script_load(settings)
    local sh = obs.obs_get_signal_handler()
    obs.signal_handler_connect(sh, "source_create", source_def.id)
end

-- OBS Studio will call this function when the script is unloaded
function script_unload()
    obs.signal_handler_disconnect(obs.obs_get_signal_handler(), "source_create", source_def.id)
end

obs.obs_register_source(source_def)