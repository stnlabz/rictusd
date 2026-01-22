import sys
import json
import requests
import cv2
from hailo_platform import VDevice, HEF, HailoSchedulingAlgorithm

# Configuration
HEF_PATH = "models/yolov8m.hef"
RICTUS_API = "http://localhost:8080/api/threat"

def run_vision_node():
    # 1. Initialize Hailo VDevice
    params = VDevice.create_params()
    params.scheduling_algorithm = HailoSchedulingAlgorithm.ROUND_ROBIN
    
    with VDevice(params) as target:
        hef = HEF(HEF_PATH)
        
        # 2. Configure Inference Model
        with target.create_infer_model(HEF_PATH) as infer_model:
            print("🚀 Hailo-8 Vision Node Active. Monitoring stream...", file=sys.stderr)
            
            # Start camera capture
            cap = cv2.VideoCapture(0)
            
            while cap.isOpened():
                ret, frame = cap.read()
                if not ret: break

                # 3. NPU Inference Processing
                # Note: In a live environment, you would use the Hailo VStreams 
                # to feed 'frame' into 'infer_model.run()'
                
                # Mock detection for logic verification
                detection = {"label": "person", "confidence": 0.96}
                
                # 4. Trigger the Triad Handoff
                if detection["confidence"] > 0.85:
                    try:
                        requests.post(RICTUS_API, json={
                            "source": "Hailo-NPU",
                            "label": detection["label"],
                            "confidence": detection["confidence"]
                        }, timeout=0.1)
                    except Exception:
                        pass # Engine busy or restarting

    cap.release()

if __name__ == "__main__":
    run_vision_node()
