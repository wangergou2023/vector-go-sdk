#!/usr/bin/env python3

"""
pip install opencv-python
pip install mediapipe
"""

import os
import sys
import cv2
import mediapipe as mp

if len(sys.argv)>0:
    fileName = sys.argv[1]

mp_hands = mp.solutions.hands
mp_drawing = mp.solutions.drawing_utils
mp_drawing_styles = mp.solutions.drawing_styles
debug = False

with mp_hands.Hands(
    static_image_mode=True,
    max_num_hands=2,
    min_detection_confidence=0.5) as hands:
        # Read an image, flip it around y-axis for correct handedness output (see
        # above).
        image = cv2.flip(cv2.imread(fileName), 1)
        # Convert the BGR image to RGB before processing.
        results = hands.process(cv2.cvtColor(image, cv2.COLOR_BGR2RGB))

        # Print handedness and draw hand landmarks on the image.
        if (debug) : 
            print('Handedness:', results.multi_handedness)
        if results.multi_hand_landmarks:
            image_height, image_width, _ = image.shape
            annotated_image = image.copy()
            for hand_landmarks in results.multi_hand_landmarks:                
                  if (debug) : 
                    print('hand_landmarks:', hand_landmarks)
                    print(
                          f'Index finger tip coordinates: (',
                          f'{hand_landmarks.landmark[mp_hands.HandLandmark.INDEX_FINGER_TIP].x * image_width}, '
                          f'{hand_landmarks.landmark[mp_hands.HandLandmark.INDEX_FINGER_TIP].y * image_height})'
                    )
                  mp_drawing.draw_landmarks(
                      annotated_image,
                      hand_landmarks,
                      mp_hands.HAND_CONNECTIONS,
                      mp_drawing_styles.get_default_hand_landmarks_style(),
                      mp_drawing_styles.get_default_hand_connections_style())
            
            saveFname = fileName.replace(".jpg", "_annotated.jpg")
            cv2.imwrite(saveFname, cv2.flip(annotated_image, 1))
            
            handList = []
            upCount = 0 
            for handLms in results.multi_hand_landmarks:
                for idx, lm in enumerate(handLms.landmark):
                    h, w, c = image.shape
                    cx, cy = int(lm.x * w), int(lm.y * h)
                    handList.append((cx, cy))
                
            finger_Coord = [(8, 5), (12, 9), (16, 13), (20, 17)]
            thumb_Coord = (4,1)
            
            i = 0
            for coordinate in finger_Coord:
                #Calculate the minimum finger length we are allowing to call it "raised"
                end = coordinate[1]
                start = end+1
                minFingerLen = (handList[end][1] - handList[start][1])*80/100
                fingerLen = handList[start][1] - handList[coordinate[0]][1]
                if (debug) :
                    print("Start("+str(start)+"): "+str(handList[start][1])+",end ("+str(end)+"):"+str(handList[end][1])) 
                    print("Min. finger length for finger "+str(i)+": "+str(minFingerLen)+", found:"+str(fingerLen))
                if minFingerLen>0 and fingerLen>minFingerLen:
                    if (debug) :
                        print("   ->"+str(i)+" finger up: "+str(handList[start][1])+"<"+str(handList[coordinate[0]][1])+", fingerlen = "+str(fingerLen)+" (min:"+str(minFingerLen)+")")
                    upCount += 1
                i = i+1
            
            #Thumb processing
            end = thumb_Coord[1]
            start = end+1
            
            handIndex = results.multi_hand_landmarks.index(hand_landmarks)
            handLabel = results.multi_handedness[handIndex].classification[0].label

            #TODO: We have also to understand whether the hand is reversed
            handFacing = "reverse"
            if (handList[17][0]>handList[5][0]):
                handFacing = "front"
                
            
            if handLabel == "Left":
                minThumbLen = (handList[end][0] - handList[start][0])*80/100
                if handFacing == "front":
                    thumbLen = handList[4][0] - handList[2][0]
                else:
                    thumbLen = handList[2][0] - handList[4][0]
            else:  
                minThumbLen = (handList[start][0]-handList[end][0])*80/100
                if handFacing == "front":
                    thumbLen = handList[2][0] - handList[4][0]
                else: 
                    thumbLen = handList[4][0] - handList[2][0]
                    
            if (debug) :
                print("Hand: "+handLabel+" facing: "+handFacing)
                print("Start("+str(start)+"): "+str(handList[start][0])+",end ("+str(end)+"):"+str(handList[end][0])) 
                print("Min. thumb size: "+str(minThumbLen)+" found:"+str(thumbLen))

            if thumbLen>minThumbLen:
                if (debug) :
                    print("Thumb up!")
                upCount += 1
            
            print(str(upCount))