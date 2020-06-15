package com.google.fcmexternalprobertarget;

import androidx.appcompat.app.AppCompatActivity;
import androidx.localbroadcastmanager.content.LocalBroadcastManager;

import android.content.IntentFilter;
import android.os.Bundle;
import android.widget.TextView;

/**
 * App that receives messages from Firebase Cloud Messaging and logs reception
 */
public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        createBroadcastReceiver();
    }

    private void createBroadcastReceiver() {
        TextView logText = findViewById(R.id.logTextView);
        LogReceiver logReceiver = new LogReceiver(logText);
        IntentFilter logFilter = new IntentFilter();

        LocalBroadcastManager.getInstance(this).registerReceiver(logReceiver, logFilter);
    }
}