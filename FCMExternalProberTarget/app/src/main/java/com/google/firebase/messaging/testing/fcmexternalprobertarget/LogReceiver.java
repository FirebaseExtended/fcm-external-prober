/*
 * Copyright 2020 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.google.firebase.messaging.testing.fcmexternalprobertarget;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.util.Log;
import android.widget.TextView;

import androidx.annotation.NonNull;

/**
 * Receives and handles broadcasts indicating an addition to the in-app logging output
 */
public class LogReceiver extends BroadcastReceiver {

    private TextView logTextView;

    /**
     * Initialize with a given UI element to modify
     * @param logTextView UI element to be modified upon receiving an appropriate broadcast
     */
    public LogReceiver(TextView logTextView) {
        this.logTextView = logTextView;
    }

    @Override
    public void onReceive(Context context, Intent intent) {
        String logText = intent.getStringExtra("logText");
        if (logText == null) {
            updateUI("Log Error: No log text supplied");
        }
        else {
            updateUI(logText);
        }
    }

    private void updateUI(@NonNull String logText) {
        int textViewLength = logTextView.length();
        logTextView.append(logText + "\n", textViewLength, textViewLength + logText.length() + 1);
    }
}
