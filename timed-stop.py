import obspython as obs
import threading
import time

# Configuration
check_interval = 60 * 60  # Check every hour
response_time = 60  # Time to respond in seconds

# Global variables
stop_recording_event = threading.Event()
is_recording = False

def remind_user():
    global is_recording
    while not stop_recording_event.is_set():
        time.sleep(check_interval)
        if is_recording:
            response = obs.obs_frontend_push_ui_translation(obs.obs_module_t())
            button = obs.obs_frontend_add_tools_menu_qaction(obs.obs_module_text("Are you still recording?"))
            button.triggered.connect(stop_recording)
            response_time_counter = 0
            while response_time_counter < response_time:
                if not is_recording:
                    break
                time.sleep(1)
                response_time_counter += 1
            if is_recording:
                stop_recording()

def start_recording(props, prop):
    global is_recording
    if not is_recording:
        obs.obs_frontend_recording_start()
        is_recording = True

def stop_recording():
    global is_recording
    if is_recording:
        obs.obs_frontend_recording_stop()
        is_recording = False

def script_description():
    return "Script to check if you are still recording and stop recording if there is no response within the set time."

def script_properties():
    props = obs.obs_properties_create()
    obs.obs_properties_add_button(props, "start_recording", "Start Recording", start_recording)
    obs.obs_properties_add_button(props, "stop_recording", "Stop Recording", stop_recording)
    return props

def script_load(settings):
    global stop_recording_event
    stop_recording_event.clear()
    remind_thread = threading.Thread(target=remind_user)
    remind_thread.start()

def script_unload():
    global stop_recording_event
    stop_recording_event.set()