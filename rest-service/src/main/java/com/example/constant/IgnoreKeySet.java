/**
 * Copyright © 2014-2017 TransfarPay.All Rights Reserved.
 */
package com.example.constant;

import java.util.HashSet;
import java.util.Set;

public class IgnoreKeySet {

	public static final Set<String> ignoreKeySet = new HashSet();
	static {
		ignoreKeySet.add("dog_ak");
		ignoreKeySet.add("appid");
		ignoreKeySet.add("dog_sk");
		ignoreKeySet.add("tf_sign");
		ignoreKeySet.add("dog_key");
		ignoreKeySet.add("tf_timestamp");
		ignoreKeySet.add("dog_value");
	}
}
