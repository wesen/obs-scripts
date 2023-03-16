obs = obslua
timers_active = false
callback_added = false

duration_s = 5

-- Function to stop the recording
function stop_recording()
    obs.script_log(obs.LOG_INFO, "Stopping recording (will show confirmation dialog)")
    show_warning()
    obs.obs_frontend_recording_stop()
    stop_timers()
    start_timers()
end

-- Function to start the timer
function start_timers()
    if not timers_active then
        local current_time = obs.os_gettime_ns()
        timer_end = current_time + duration_s * 1000000000 -- 1 hour in nanoseconds
        obs.timer_add(stop_recording, duration_s * 1000) -- 1 hour in milliseconds
        obs.timer_add(log_remaining_time, 1000) -- Log remaining time every second
        timers_active = true
    end
end

function show_warning()
    local scenes = obs.obs_frontend_get_scenes()
    if scenes ~= nil then
        local scene = obs.obs_frontend_get_current_scene()
        local text_source = obs.obs_get_source_by_name(source_name)
        if text_source == nil then
            local text_settings = obs.obs_data_create()
            obs.obs_data_set_string(text_settings, "text", "Click YES in the script settings within 1 minute to continue recording.")
            
            text_source = obs.obs_source_create_private("text_gdiplus", source_name, text_settings)
            obs.obs_data_release(text_settings)
            
            local scene_item = obs.obs_scene_add(scene, text_source)
            obs.obs_sceneitem_set_visible(scene_item, true)
        end
        obs.obs_source_release(text_source)
        obs.obs_scene_release(scene)
    end
end

function stop_timers() 
    if timers_active then
        obs.timer_remove(stop_recording)
        obs.timer_remove(log_remaining_time)
        timers_active = false
    end
end

-- Function to log the remaining time
function log_remaining_time()
    local current_time = obs.os_gettime_ns()
    local remaining_time = (timer_end - current_time) / 1000000000 -- Convert from nanoseconds to seconds
    local minutes = math.floor(remaining_time / 60)
    local seconds = math.floor(remaining_time % 60)
    obs.script_log(obs.LOG_INFO, string.format("Remaining time: %02d:%02d", minutes, seconds))
end

-- Function called when the recording starts
function on_recording_started(event)
    if event == obs.OBS_FRONTEND_EVENT_RECORDING_STARTED then
        obs.script_log(obs.LOG_INFO, "Recording started")
        start_timers()
    else if event == obs.OBS_FRONTEND_EVENT_RECORDING_STOPPED then
        obs.script_log(obs.LOG_INFO, "Recording stopped")
        stop_timers()
    end
    end
end

-- OBS Studio will call this function when the script is loaded
function script_load(settings)
    if not callback_added then
        obs.script_log(obs.LOG_INFO, "Adding callback")
        obs.obs_frontend_add_event_callback(on_recording_started)
        callback_added = true
    end
end

-- OBS Studio will call this function when the script is unloaded
function script_unload()
    stop_timers()
    if callback_added then
        -- obs.script_log(obs.LOG_INFO, "Removing callback")
        -- obs.obs_frontend_remove_event_callback(on_recording_started)
        callback_added = false  
    end
end
