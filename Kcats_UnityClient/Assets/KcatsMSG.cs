using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class KcatsMSG
{
    private string Sender;
    private string Receiver;
    private string Text;

    public void Set(string sr, string rc, string text)
    {
        Sender = sr;
        Receiver = rc;
        Text = text;
    }

    public string GetMainText()
    {
        return Text;
    }

    public string msg2Text()
    {
        string str = "\n  " + "<color=blue>" + Sender + "</color>: " + Text;
        return str;
    }

    public string ClientRegister()
    {
        return "";
    }

    public bool IsServerCall()
    {
        if (Sender == "NIL" || Sender == "Server")
            return true;
        return false;
    }
}
