using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class ServerCall
{
    public string[] DispatchServerCall(string text)
    {
        string[] all = text.Split('\r');
        string[] slice = all[0].Split(' ');
        switch (slice[0])
        {
            case "C": // 用于同步用户列表
                return slice;
            case "E":
                return slice;
        }
        return null;
    }
}
